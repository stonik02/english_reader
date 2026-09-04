package listcatalog

import (
	"context"

	domain "github.com/deniskrylov/english-reader/backend/internal/domain/library"
)

type UseCase interface {
	Execute(context.Context, string, int) (domain.Page[domain.Book], error)
}
