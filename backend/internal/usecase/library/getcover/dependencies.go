package getcover

import (
	"context"

	domain "github.com/deniskrylov/english-reader/backend/internal/domain/library"
)

type Books interface {
	Get(context.Context, string) (domain.Book, error)
}

type Covers interface {
	Read(string) ([]byte, error)
}
