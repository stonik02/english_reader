package getbook

import (
	"context"

	domain "github.com/deniskrylov/english-reader/backend/internal/domain/library"
)

type Books interface {
	Get(context.Context, string) (domain.Book, error)
}
