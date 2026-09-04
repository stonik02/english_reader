package translatetext

import (
	"context"
	"time"

	domain "github.com/deniskrylov/english-reader/backend/internal/domain/dictionary"
)

type TranslationCache interface {
	CachedTranslation(context.Context, string, string) (domain.CachedTranslation, bool, error)
	PutTranslation(context.Context, string, string, string, time.Duration) error
}

type TranslationProvider interface {
	Translate(context.Context, string) (string, error)
}

type ReaderPort interface {
	ChapterPlainText(context.Context, string, string) (string, error)
}
