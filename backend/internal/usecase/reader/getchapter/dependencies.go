package getchapter

import (
	"context"
	domain "github.com/deniskrylov/english-reader/backend/internal/domain/reader"
)

type ReaderRepository interface {
	EnsureReadyBook(context.Context, string, string) error
	Chapter(context.Context, string, string, string) (domain.Chapter, error)
}
