package reader

import (
	"context"
	"errors"
	"fmt"

	domain "github.com/deniskrylov/english-reader/backend/internal/domain/reader"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) EnsureReadyBook(ctx context.Context, userID, bookID string) error {
	var status string
	err := r.pool.QueryRow(ctx, `SELECT status FROM books WHERE id=$1`, bookID).Scan(&status)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ErrNotFound
	}
	if err != nil {
		return err
	}
	if status != "ready" {
		return domain.ErrNotReady
	}
	_, err = r.pool.Exec(ctx, `INSERT INTO user_books (user_id,book_id,added_via) VALUES ($1,$2,'first_read') ON CONFLICT (user_id,book_id) DO NOTHING`, userID, bookID)
	return err
}

func (r *Repository) FirstChapter(ctx context.Context, bookID string) (domain.Chapter, error) {
	return r.chapter(ctx, `SELECT id,href,sequence,start_cfi,end_cfi,sanitized_html,(SELECT COUNT(*)::int FROM book_chapters WHERE book_id=$1::uuid) FROM book_chapters WHERE book_id=$1::uuid ORDER BY sequence LIMIT 1`, bookID)
}
func (r *Repository) Chapter(ctx context.Context, bookID, chapterID, href string) (domain.Chapter, error) {
	return r.chapter(ctx, `SELECT id,href,sequence,start_cfi,end_cfi,sanitized_html,(SELECT COUNT(*)::int FROM book_chapters WHERE book_id=$1::uuid) FROM book_chapters WHERE book_id=$1::uuid AND (NULLIF($2, '')::uuid IS NULL OR id=NULLIF($2, '')::uuid) AND ($3='' OR href=$3) LIMIT 1`, bookID, chapterID, href)
}
func (r *Repository) Adjacent(ctx context.Context, bookID, chapterID string, direction int) (domain.Chapter, error) {
	return r.chapter(ctx, `SELECT id,href,sequence,start_cfi,end_cfi,sanitized_html,(SELECT COUNT(*)::int FROM book_chapters WHERE book_id=$1::uuid) FROM book_chapters WHERE book_id=$1::uuid AND sequence=(SELECT sequence+$3 FROM book_chapters WHERE id=$2 AND book_id=$1::uuid)`, bookID, chapterID, direction)
}
func (r *Repository) Progress(ctx context.Context, userID, bookID string) (domain.Progress, error) {
	var value domain.Progress
	err := r.pool.QueryRow(ctx, `SELECT COALESCE(chapter_id::text,''),epub_cfi,progress_percent,revision,updated_at FROM reading_progress WHERE user_id=$1 AND book_id=$2`, userID, bookID).Scan(&value.ChapterID, &value.EPUBCFI, &value.ProgressPercent, &value.Revision, &value.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Progress{}, domain.ErrNotFound
	}
	return value, err
}
func (r *Repository) SaveProgress(ctx context.Context, userID, bookID string, value domain.Progress) (domain.Progress, error) {
	err := r.pool.QueryRow(ctx, `INSERT INTO reading_progress (user_id,book_id,chapter_id,epub_cfi,progress_percent,revision) VALUES ($1,$2,$3,$4,$5,$6) ON CONFLICT (user_id,book_id) DO UPDATE SET chapter_id=EXCLUDED.chapter_id,epub_cfi=EXCLUDED.epub_cfi,progress_percent=EXCLUDED.progress_percent,revision=EXCLUDED.revision,updated_at=NOW() WHERE reading_progress.revision < EXCLUDED.revision RETURNING chapter_id,epub_cfi,progress_percent,revision,updated_at`, userID, bookID, value.ChapterID, value.EPUBCFI, value.ProgressPercent, value.Revision).Scan(&value.ChapterID, &value.EPUBCFI, &value.ProgressPercent, &value.Revision, &value.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return r.Progress(ctx, userID, bookID)
	}
	return value, err
}
func (r *Repository) Settings(ctx context.Context, userID string) (domain.Settings, error) {
	var value domain.Settings
	err := r.pool.QueryRow(ctx, `INSERT INTO user_reader_settings (user_id) VALUES ($1) ON CONFLICT (user_id) DO UPDATE SET user_id=EXCLUDED.user_id RETURNING font_scale,theme,line_height,highlight_color`, userID).Scan(&value.FontScale, &value.Theme, &value.LineHeight, &value.HighlightColor)
	return value, err
}
func (r *Repository) UpdateSettings(ctx context.Context, userID string, value domain.Settings) (domain.Settings, error) {
	err := r.pool.QueryRow(ctx, `INSERT INTO user_reader_settings (user_id,font_scale,theme,line_height,highlight_color) VALUES ($1,$2,$3,$4,$5) ON CONFLICT (user_id) DO UPDATE SET font_scale=EXCLUDED.font_scale,theme=EXCLUDED.theme,line_height=EXCLUDED.line_height,highlight_color=EXCLUDED.highlight_color,updated_at=NOW() RETURNING font_scale,theme,line_height,highlight_color`, userID, value.FontScale, value.Theme, value.LineHeight, value.HighlightColor).Scan(&value.FontScale, &value.Theme, &value.LineHeight, &value.HighlightColor)
	return value, err
}
func (r *Repository) chapter(ctx context.Context, query string, args ...any) (domain.Chapter, error) {
	var value domain.Chapter
	err := r.pool.QueryRow(ctx, query, args...).Scan(&value.ID, &value.Href, &value.Sequence, &value.StartCFI, &value.EndCFI, &value.SanitizedHTML, &value.TotalChapters)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Chapter{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.Chapter{}, fmt.Errorf("read chapter: %w", err)
	}
	return value, nil
}
