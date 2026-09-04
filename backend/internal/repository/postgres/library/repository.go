package library

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	domain "github.com/deniskrylov/english-reader/backend/internal/domain/library"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct{ pool *pgxpool.Pool }

func New(pool *pgxpool.Pool) *Repository { return &Repository{pool: pool} }

func (r *Repository) FindBySHA256(ctx context.Context, hash []byte) (domain.Book, error) {
	return r.findBook(ctx, `SELECT id,title,author,status,COALESCE(cover_file_path,''),uploaded_by_user_id,created_at FROM books WHERE content_sha256=$1`, hash)
}

func (r *Repository) Create(ctx context.Context, id, userID, title, path string, hash []byte) (domain.Book, error) {
	_, err := r.pool.Exec(ctx, `INSERT INTO books (id,uploaded_by_user_id,title,format,status,content_sha256,source_file_path) VALUES ($1,$2,$3,'epub','processing',$4,$5)`, id, userID, title, hash, path)
	if err != nil {
		return domain.Book{}, fmt.Errorf("create book: %w", err)
	}
	_, err = r.pool.Exec(ctx, `INSERT INTO ingestion_jobs (id,book_id,state) VALUES ($1,$2,'pending')`, uuid.NewString(), id)
	if err != nil {
		return domain.Book{}, fmt.Errorf("create ingestion job: %w", err)
	}
	return r.Get(ctx, id)
}

func (r *Repository) ClaimNextIngestion(ctx context.Context) (domain.IngestionJob, error) {
	var job domain.IngestionJob
	err := r.pool.QueryRow(ctx, `
		WITH next AS (
			SELECT book_id FROM ingestion_jobs
			WHERE state IN ('pending', 'processing')
			ORDER BY created_at
			FOR UPDATE SKIP LOCKED
			LIMIT 1
		)
		UPDATE ingestion_jobs job
		SET state='processing', attempt=attempt+1, locked_at=NOW(), updated_at=NOW()
		FROM next JOIN books book ON book.id=next.book_id
		WHERE job.book_id=next.book_id
		RETURNING job.book_id, book.source_file_path`).Scan(&job.BookID, &job.SourceFilePath)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.IngestionJob{}, domain.ErrNotFound
	}
	return job, err
}

func (r *Repository) CompleteIngestion(ctx context.Context, bookID, title, author, coverPath string, chapters []domain.Chapter) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if _, err = tx.Exec(ctx, `DELETE FROM book_chapters WHERE book_id=$1`, bookID); err != nil {
		return err
	}
	for _, chapter := range chapters {
		if _, err = tx.Exec(ctx, `INSERT INTO book_chapters (id,book_id,sequence,href,start_cfi,end_cfi,sanitized_html,plain_text,word_count) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`, uuid.NewString(), bookID, chapter.Sequence, chapter.Href, chapter.StartCFI, chapter.EndCFI, chapter.SanitizedHTML, chapter.PlainText, len(strings.Fields(chapter.PlainText))); err != nil {
			return err
		}
	}
	result, err := tx.Exec(ctx, `UPDATE books SET title=$2,author=$3,cover_file_path=NULLIF($4,''),status='ready',failure_code=NULL,updated_at=NOW() WHERE id=$1`, bookID, title, author, coverPath)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	if _, err = tx.Exec(ctx, `UPDATE ingestion_jobs SET state='completed',finished_at=NOW(),updated_at=NOW(),error_detail=NULL WHERE book_id=$1`, bookID); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *Repository) FailIngestion(ctx context.Context, bookID, code string) error {
	_, err := r.pool.Exec(ctx, `UPDATE books SET status='failed',failure_code=$2,updated_at=NOW() WHERE id=$1`, bookID, code)
	if err != nil {
		return err
	}
	_, err = r.pool.Exec(ctx, `UPDATE ingestion_jobs SET state='failed',finished_at=NOW(),updated_at=NOW(),error_detail=$2 WHERE book_id=$1`, bookID, code)
	return err
}

// RequeueIngestion schedules an existing book for parsing from its stored
// original EPUB. It intentionally keeps the original file and user-book links.
func (r *Repository) RequeueIngestion(ctx context.Context, bookID string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	result, err := tx.Exec(ctx, `UPDATE books SET status='processing',failure_code=NULL,updated_at=NOW() WHERE id=$1`, bookID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	result, err = tx.Exec(ctx, `UPDATE ingestion_jobs SET state='pending',locked_at=NULL,finished_at=NULL,error_detail=NULL,updated_at=NOW() WHERE book_id=$1`, bookID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("book %s has no ingestion job", bookID)
	}
	return tx.Commit(ctx)
}

// RequeueAllIngestions schedules all stored EPUBs for reparsing. It is an
// operator action used after parser improvements, not part of the public API.
func (r *Repository) RequeueAllIngestions(ctx context.Context) (int64, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)
	result, err := tx.Exec(ctx, `UPDATE books SET status='processing',failure_code=NULL,updated_at=NOW()`)
	if err != nil {
		return 0, err
	}
	if _, err = tx.Exec(ctx, `UPDATE ingestion_jobs SET state='pending',locked_at=NULL,finished_at=NULL,error_detail=NULL,updated_at=NOW()`); err != nil {
		return 0, err
	}
	if err = tx.Commit(ctx); err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

func (r *Repository) Get(ctx context.Context, bookID string) (domain.Book, error) {
	return r.findBook(ctx, `SELECT id,title,author,status,COALESCE(cover_file_path,''),uploaded_by_user_id,created_at FROM books WHERE id=$1`, bookID)
}

func (r *Repository) ListCatalog(ctx context.Context, cursor string, limit int) (domain.Page[domain.Book], error) {
	return r.listBooks(ctx, `SELECT id,title,author,status,COALESCE(cover_file_path,''),uploaded_by_user_id,created_at FROM books WHERE (NULLIF($1, '') IS NULL OR id < NULLIF($1, '')::uuid) ORDER BY id DESC LIMIT $2`, cursor, limit)
}

func (r *Repository) Add(ctx context.Context, userID, bookID, via string) (domain.UserBook, error) {
	_, err := r.pool.Exec(ctx, `INSERT INTO user_books (user_id,book_id,added_via) VALUES ($1,$2,$3) ON CONFLICT (user_id,book_id) DO NOTHING`, userID, bookID, via)
	if err != nil {
		return domain.UserBook{}, fmt.Errorf("add book: %w", err)
	}
	return r.userBook(ctx, userID, bookID)
}

func (r *Repository) ListMine(ctx context.Context, userID, cursor string, limit int) (domain.Page[domain.UserBook], error) {
	rows, err := r.pool.Query(ctx, `SELECT b.id,b.title,b.author,b.status,COALESCE(b.cover_file_path,''),b.uploaded_by_user_id,b.created_at,ub.added_at,ub.added_via,COALESCE(rp.progress_percent,0) FROM user_books ub JOIN books b ON b.id=ub.book_id LEFT JOIN reading_progress rp ON rp.user_id=ub.user_id AND rp.book_id=ub.book_id WHERE ub.user_id=$1 AND (NULLIF($2, '') IS NULL OR b.id < NULLIF($2, '')::uuid) ORDER BY b.id DESC LIMIT $3`, userID, cursor, limit+1)
	if err != nil {
		return domain.Page[domain.UserBook]{}, err
	}
	defer rows.Close()
	page := domain.Page[domain.UserBook]{Items: make([]domain.UserBook, 0, limit)}
	for rows.Next() {
		var item domain.UserBook
		if err := rows.Scan(&item.Book.ID, &item.Book.Title, &item.Book.Author, &item.Book.Status, &item.Book.CoverPath, &item.Book.UploadedByUserID, &item.Book.CreatedAt, &item.AddedAt, &item.AddedVia, &item.ProgressPercent); err != nil {
			return page, err
		}
		item.Book.CoverURL = coverURL(item.Book.ID, item.Book.CoverPath)
		page.Items = append(page.Items, item)
	}
	if len(page.Items) > limit {
		page.NextCursor = page.Items[limit-1].Book.ID
		page.Items = page.Items[:limit]
	}
	return page, rows.Err()
}

func (r *Repository) Remove(ctx context.Context, userID, bookID string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM user_books WHERE user_id=$1 AND book_id=$2`, userID, bookID)
	return err
}

func (r *Repository) Delete(ctx context.Context, bookID string) (domain.StoredBookFiles, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return domain.StoredBookFiles{}, err
	}
	defer tx.Rollback(ctx)
	var files domain.StoredBookFiles
	err = tx.QueryRow(ctx, `DELETE FROM books WHERE id=$1 RETURNING source_file_path,COALESCE(cover_file_path,'')`, bookID).Scan(&files.SourcePath, &files.CoverPath)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.StoredBookFiles{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.StoredBookFiles{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return domain.StoredBookFiles{}, err
	}
	return files, nil
}

func (r *Repository) findBook(ctx context.Context, query string, value any) (domain.Book, error) {
	var book domain.Book
	err := r.pool.QueryRow(ctx, query, value).Scan(&book.ID, &book.Title, &book.Author, &book.Status, &book.CoverPath, &book.UploadedByUserID, &book.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Book{}, domain.ErrNotFound
	}
	book.CoverURL = coverURL(book.ID, book.CoverPath)
	return book, err
}

func (r *Repository) listBooks(ctx context.Context, query, cursor string, limit int) (domain.Page[domain.Book], error) {
	rows, err := r.pool.Query(ctx, query, cursor, limit+1)
	if err != nil {
		return domain.Page[domain.Book]{}, err
	}
	defer rows.Close()
	page := domain.Page[domain.Book]{Items: make([]domain.Book, 0, limit)}
	for rows.Next() {
		var book domain.Book
		if err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Status, &book.CoverPath, &book.UploadedByUserID, &book.CreatedAt); err != nil {
			return page, err
		}
		book.CoverURL = coverURL(book.ID, book.CoverPath)
		page.Items = append(page.Items, book)
	}
	if len(page.Items) > limit {
		page.NextCursor = page.Items[limit-1].ID
		page.Items = page.Items[:limit]
	}
	return page, rows.Err()
}

func (r *Repository) userBook(ctx context.Context, userID, bookID string) (domain.UserBook, error) {
	rows, err := r.pool.Query(ctx, `SELECT b.id,b.title,b.author,b.status,COALESCE(b.cover_file_path,''),b.uploaded_by_user_id,b.created_at,ub.added_at,ub.added_via,COALESCE(rp.progress_percent,0) FROM user_books ub JOIN books b ON b.id=ub.book_id LEFT JOIN reading_progress rp ON rp.user_id=ub.user_id AND rp.book_id=ub.book_id WHERE ub.user_id=$1 AND ub.book_id=$2`, userID, bookID)
	if err != nil {
		return domain.UserBook{}, err
	}
	defer rows.Close()
	if !rows.Next() {
		return domain.UserBook{}, domain.ErrNotFound
	}
	var item domain.UserBook
	err = rows.Scan(&item.Book.ID, &item.Book.Title, &item.Book.Author, &item.Book.Status, &item.Book.CoverPath, &item.Book.UploadedByUserID, &item.Book.CreatedAt, &item.AddedAt, &item.AddedVia, &item.ProgressPercent)
	item.Book.CoverURL = coverURL(item.Book.ID, item.Book.CoverPath)
	return item, err
}

func coverURL(bookID, coverPath string) string {
	if coverPath == "" {
		return ""
	}
	return "/api/v1/library/books/" + bookID + "/cover"
}

func SHA256String(hash [32]byte) string { return hex.EncodeToString(hash[:]) }
