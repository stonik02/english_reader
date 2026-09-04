package getbook

import (
	"context"

	domain "github.com/deniskrylov/english-reader/backend/internal/domain/library"
)

type UseCase interface {
	Execute(context.Context, string) (domain.Book, error)
}
