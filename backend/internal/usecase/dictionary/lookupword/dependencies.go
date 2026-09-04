package lookupword

import (
	"context"
	"time"

	domain "github.com/deniskrylov/english-reader/backend/internal/domain/dictionary"
)

type WordNormalizer interface{ Normalize(string) (string, error) }
type Morphology interface{ Candidates(string) []string }
type DictionaryRepository interface {
	Lookup(context.Context, string) (domain.LookupResult, error)
	CachedTranslation(context.Context, string, string) (domain.CachedTranslation, bool, error)
	PutTranslation(context.Context, string, string, string, time.Duration) error
}
type TranslationProvider interface {
	Translate(context.Context, string) (string, error)
}
type ReaderPort interface {
	ChapterPlainText(context.Context, string, string) (string, error)
}

type VocabularyReader interface {
	IsSaved(context.Context, string, int64) (bool, error)
}
