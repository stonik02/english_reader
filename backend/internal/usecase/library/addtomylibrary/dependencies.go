package addtomylibrary

import (
	"context"

	domain "github.com/deniskrylov/english-reader/backend/internal/domain/library"
)

type Books interface {
	Get(context.Context, string) (domain.Book, error)
	Add(context.Context, string, string, string) (domain.UserBook, error)
}
