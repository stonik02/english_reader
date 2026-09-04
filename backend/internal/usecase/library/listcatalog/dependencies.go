package listcatalog

import (
	"context"

	domain "github.com/deniskrylov/english-reader/backend/internal/domain/library"
)

type Books interface {
	ListCatalog(context.Context, string, int) (domain.Page[domain.Book], error)
}
