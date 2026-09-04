package vocabulary

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	domain "github.com/deniskrylov/english-reader/backend/internal/domain/vocabulary"
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Save(ctx context.Context, userID string, request domain.SaveRequest) (domain.Entry, bool, error) {
	if err := r.validateSave(ctx, request); err != nil {
		return domain.Entry{}, false, err
	}

	entryID := uuid.NewString()
	inserted := false
	err := r.pool.QueryRow(ctx, `
		INSERT INTO vocabulary_entries (id, user_id, lemma_id, chosen_sense_id, source_form)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id, lemma_id) DO NOTHING
		RETURNING true`, entryID, userID, request.LemmaID, request.ChosenSenseID, request.SourceForm).Scan(&inserted)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return domain.Entry{}, false, fmt.Errorf("insert vocabulary entry: %w", err)
	}

	entry, err := r.entryByLemma(ctx, userID, request.LemmaID)
	if err != nil {
		return domain.Entry{}, false, err
	}

	return entry, !inserted, nil
}

func (r *Repository) List(ctx context.Context, userID string, request domain.ListRequest) ([]domain.Entry, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT ve.id, ve.lemma_id, dl.lemma, ve.source_form, ve.created_at, ve.updated_at,
			   ds.id, ds.part_of_speech, ds.translations, COALESCE(ds.example_en, ''), COALESCE(ds.example_ru, '')
		FROM vocabulary_entries ve
		JOIN dictionary_lemmas dl ON dl.id = ve.lemma_id
		LEFT JOIN LATERAL (
			SELECT id, part_of_speech, translations, example_en, example_ru
			FROM dictionary_senses
			WHERE lemma_id = ve.lemma_id
			  AND (ve.chosen_sense_id IS NULL OR id = ve.chosen_sense_id)
			ORDER BY position ASC
			LIMIT 1
		) ds ON true
		WHERE ve.user_id = $1
		  AND ($2 = '' OR dl.lemma ILIKE '%' || $2 || '%')
		  AND ($3::timestamptz IS NULL OR (ve.created_at, ve.id) < ($3::timestamptz, $4::uuid))
		ORDER BY ve.created_at DESC, ve.id DESC
		LIMIT $5`, userID, request.Query, cursorTime(request.Cursor), cursorID(request.Cursor), request.Limit)
	if err != nil {
		return nil, fmt.Errorf("list vocabulary entries: %w", err)
	}
	defer rows.Close()

	entries := make([]domain.Entry, 0)
	for rows.Next() {
		entry, err := scanEntry(rows)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate vocabulary entries: %w", err)
	}

	return entries, nil
}

func (r *Repository) DeleteByID(ctx context.Context, userID, entryID string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM vocabulary_entries WHERE user_id = $1 AND id = $2`, userID, entryID)
	return err
}

func (r *Repository) DeleteByLemmaID(ctx context.Context, userID string, lemmaID int64) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM vocabulary_entries WHERE user_id = $1 AND lemma_id = $2`, userID, lemmaID)
	return err
}

func (r *Repository) HighlightLemmas(ctx context.Context, userID string, limit int) ([]domain.HighlightLemma, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT dl.id, dl.lemma
		FROM vocabulary_entries ve
		JOIN dictionary_lemmas dl ON dl.id = ve.lemma_id
		WHERE ve.user_id = $1
		ORDER BY ve.created_at DESC, ve.id DESC
		LIMIT $2`, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("list vocabulary highlight lemmas: %w", err)
	}
	defer rows.Close()

	lemmas := make([]domain.HighlightLemma, 0)
	for rows.Next() {
		var lemma domain.HighlightLemma
		if err := rows.Scan(&lemma.ID, &lemma.Lemma); err != nil {
			return nil, fmt.Errorf("scan vocabulary highlight lemma: %w", err)
		}
		lemmas = append(lemmas, lemma)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate vocabulary highlight lemmas: %w", err)
	}

	return lemmas, nil
}

func (r *Repository) IsSaved(ctx context.Context, userID string, lemmaID int64) (bool, error) {
	var saved bool
	err := r.pool.QueryRow(ctx, `
		SELECT EXISTS(
			SELECT 1
			FROM vocabulary_entries
			WHERE user_id = $1 AND lemma_id = $2
		)`, userID, lemmaID).Scan(&saved)
	if err != nil {
		return false, fmt.Errorf("check saved vocabulary entry: %w", err)
	}

	return saved, nil
}

func (r *Repository) ChapterPlainText(ctx context.Context, userID, bookID, chapterID string) (string, error) {
	var status string
	err := r.pool.QueryRow(ctx, `
		SELECT b.status
		FROM user_books ub
		JOIN books b ON b.id = ub.book_id
		WHERE ub.user_id = $1 AND ub.book_id = $2`, userID, bookID).Scan(&status)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", domain.ErrNotFound
	}
	if err != nil {
		return "", fmt.Errorf("get book status for highlights: %w", err)
	}
	if status != "ready" {
		return "", domain.ErrNotReady
	}

	var plainText string
	err = r.pool.QueryRow(ctx, `
		SELECT plain_text
		FROM book_chapters
		WHERE book_id = $1 AND id = $2`, bookID, chapterID).Scan(&plainText)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", domain.ErrNotFound
	}
	if err != nil {
		return "", fmt.Errorf("get chapter plain text for highlights: %w", err)
	}

	return plainText, nil
}

func (r *Repository) validateSave(ctx context.Context, request domain.SaveRequest) error {
	var exists bool
	err := r.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM dictionary_lemmas WHERE id = $1 AND language = 'en')`, request.LemmaID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("check dictionary lemma: %w", err)
	}
	if !exists {
		return domain.ErrNotFound
	}
	if request.ChosenSenseID == nil {
		return nil
	}
	err = r.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM dictionary_senses WHERE id = $1 AND lemma_id = $2)`, *request.ChosenSenseID, request.LemmaID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("check dictionary sense: %w", err)
	}
	if !exists {
		return domain.ErrInvalidSense
	}

	return nil
}

func (r *Repository) entryByLemma(ctx context.Context, userID string, lemmaID int64) (domain.Entry, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT ve.id, ve.lemma_id, dl.lemma, ve.source_form, ve.created_at, ve.updated_at,
			   ds.id, ds.part_of_speech, ds.translations, COALESCE(ds.example_en, ''), COALESCE(ds.example_ru, '')
		FROM vocabulary_entries ve
		JOIN dictionary_lemmas dl ON dl.id = ve.lemma_id
		LEFT JOIN LATERAL (
			SELECT id, part_of_speech, translations, example_en, example_ru
			FROM dictionary_senses
			WHERE lemma_id = ve.lemma_id
			  AND (ve.chosen_sense_id IS NULL OR id = ve.chosen_sense_id)
			ORDER BY position ASC
			LIMIT 1
		) ds ON true
		WHERE ve.user_id = $1 AND ve.lemma_id = $2`, userID, lemmaID)
	entry, err := scanEntry(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Entry{}, domain.ErrNotFound
	}
	return entry, err
}

type rowScanner interface {
	Scan(...any) error
}

func scanEntry(row rowScanner) (domain.Entry, error) {
	var entry domain.Entry
	var senseID *int64
	var partOfSpeech *string
	var translations []byte
	var exampleEN *string
	var exampleRU *string
	if err := row.Scan(
		&entry.ID,
		&entry.LemmaID,
		&entry.Lemma,
		&entry.SourceForm,
		&entry.CreatedAt,
		&entry.UpdatedAt,
		&senseID,
		&partOfSpeech,
		&translations,
		&exampleEN,
		&exampleRU,
	); err != nil {
		return domain.Entry{}, err
	}
	if senseID == nil {
		return entry, nil
	}

	translationsValue := make([]string, 0)
	if err := json.Unmarshal(translations, &translationsValue); err != nil {
		return domain.Entry{}, fmt.Errorf("decode selected sense translations: %w", err)
	}
	entry.ChosenSense = &domain.Sense{
		ID:           *senseID,
		PartOfSpeech: *partOfSpeech,
		Translations: translationsValue,
		ExampleEN:    *exampleEN,
		ExampleRU:    *exampleRU,
	}

	return entry, nil
}

func cursorTime(cursor *domain.Cursor) any {
	if cursor == nil {
		return nil
	}
	return cursor.CreatedAt
}

func cursorID(cursor *domain.Cursor) any {
	if cursor == nil {
		return nil
	}
	return cursor.ID
}
