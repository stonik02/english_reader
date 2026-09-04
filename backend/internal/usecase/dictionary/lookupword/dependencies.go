package lookupword

import (
	"context"

	domain "github.com/deniskrylov/english-reader/backend/internal/domain/dictionary"
)

type WordNormalizer interface{ Normalize(string) (string, error) }
type Morphology interface{ Candidates(string) []string }
type DictionaryRepository interface {
	Lookup(context.Context, string) (domain.LookupResult, error)
}

type VocabularyReader interface {
	IsSaved(context.Context, string, int64) (bool, error)
}
