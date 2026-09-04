package saveprogress

import (
	"context"
	domain "github.com/deniskrylov/english-reader/backend/internal/domain/reader"
)

type ReaderRepository interface {
	EnsureReadyBook(context.Context, string, string) error
	Chapter(context.Context, string, string, string) (domain.Chapter, error)
	SaveProgress(context.Context, string, string, domain.Progress) (domain.Progress, error)
}
