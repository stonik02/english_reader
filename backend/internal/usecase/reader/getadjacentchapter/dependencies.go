package getadjacentchapter

import (
	"context"
	domain "github.com/deniskrylov/english-reader/backend/internal/domain/reader"
)

type ReaderRepository interface {
	EnsureReadyBook(context.Context, string, string) error
	Adjacent(context.Context, string, string, int) (domain.Chapter, error)
}
