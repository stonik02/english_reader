package dictionary

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"time"

	domain "github.com/deniskrylov/english-reader/backend/internal/domain/dictionary"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Lookup(ctx context.Context, lemma string) (domain.LookupResult, error) {
	rows, err := r.pool.Query(ctx, `
		WITH latest_lemma AS (
			SELECT id,lemma,source,source_version
			FROM dictionary_lemmas
			WHERE language='en' AND lemma=$1
			ORDER BY source_version DESC,id DESC
			LIMIT 1
		), unique_senses AS (
			SELECT DISTINCT ON (s.part_of_speech,s.translations)
				s.id,s.lemma_id,s.part_of_speech,s.translations,s.example_en,s.example_ru,s.source_url,s.attribution,s.position
			FROM dictionary_senses s
			JOIN latest_lemma l ON l.id=s.lemma_id
			ORDER BY s.part_of_speech,s.translations,
				(COALESCE(s.example_en,'')<>'' OR COALESCE(s.example_ru,'')<>'') DESC,
				s.position ASC,s.id ASC
		)
		SELECT l.id,l.lemma,l.source,l.source_version,s.id,s.part_of_speech,s.translations,COALESCE(s.example_en,''),COALESCE(s.example_ru,''),s.source_url,s.attribution
		FROM latest_lemma l JOIN unique_senses s ON s.lemma_id=l.id
		ORDER BY s.position ASC,s.id ASC`, lemma)
	if err != nil {
		return domain.LookupResult{}, err
	}
	defer rows.Close()
	result := domain.LookupResult{Lemma: lemma}
	for rows.Next() {
		var sense domain.Sense
		var translations []byte
		if err := rows.Scan(&result.LemmaID, &result.Lemma, &result.Source, &result.SourceVersion, &sense.ID, &sense.PartOfSpeech, &translations, &sense.ExampleEN, &sense.ExampleRU, &sense.SourceURL, &sense.Attribution); err != nil {
			return result, err
		}
		if err := json.Unmarshal(translations, &sense.Translations); err != nil {
			return result, err
		}
		result.Senses = append(result.Senses, sense)
	}
	if err := rows.Err(); err != nil {
		return result, err
	}
	return result, nil
}

func (r *Repository) CachedTranslation(ctx context.Context, text, modelVersion string) (domain.CachedTranslation, bool, error) {
	key := cacheKey(text, modelVersion)
	var value domain.CachedTranslation
	err := r.pool.QueryRow(ctx, `SELECT translated_text,expires_at FROM translation_cache WHERE key_hash=$1 AND expires_at>NOW()`, key).Scan(&value.Text, &value.ExpiresAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.CachedTranslation{}, false, nil
	}
	return value, err == nil, err
}

func (r *Repository) PutTranslation(ctx context.Context, text, translated, modelVersion string, ttl time.Duration) error {
	key := cacheKey(text, modelVersion)
	_, err := r.pool.Exec(ctx, `INSERT INTO translation_cache (key_hash,model_version,source_lang,target_lang,normalized_text,translated_text,expires_at) VALUES ($1,$2,'en','ru',$3,$4,NOW()+$5::interval) ON CONFLICT (key_hash) DO UPDATE SET translated_text=EXCLUDED.translated_text,expires_at=EXCLUDED.expires_at,model_version=EXCLUDED.model_version`, key, modelVersion, text, translated, ttl.String())
	return err
}

func (r *Repository) ChapterPlainText(ctx context.Context, bookID, chapterID string) (string, error) {
	var status string
	err := r.pool.QueryRow(ctx, `SELECT status FROM books WHERE id=$1`, bookID).Scan(&status)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", errors.New("book not found")
	}
	if err != nil {
		return "", err
	}
	if status != "ready" {
		return "", errors.New("book is not ready")
	}
	var plainText string
	err = r.pool.QueryRow(ctx, `SELECT plain_text FROM book_chapters WHERE book_id=$1 AND id=$2`, bookID, chapterID).Scan(&plainText)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", errors.New("chapter not found")
	}
	return plainText, err
}

func cacheKey(text, modelVersion string) []byte {
	value := sha256.Sum256([]byte("en\x00ru\x00" + modelVersion + "\x00" + text))
	return value[:]
}
