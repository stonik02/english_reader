package getreadingstate

import (
	"context"

	domain "github.com/deniskrylov/english-reader/backend/internal/domain/reader"
)

type ReaderRepository interface {
	EnsureReadyBook(context.Context, string, string) error
	FirstChapter(context.Context, string) (domain.Chapter, error)
	Chapter(context.Context, string, string, string) (domain.Chapter, error)
	Progress(context.Context, string, string) (domain.Progress, error)
	Settings(context.Context, string) (domain.Settings, error)
}
