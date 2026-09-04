package listcatalog

import (
	"context"

	domain "github.com/deniskrylov/english-reader/backend/internal/domain/library"
)

type UseCase struct {
	books Books
}

func New(books Books) *UseCase {
	return &UseCase{books: books}
}

func (u *UseCase) Execute(ctx context.Context, cursor string, limit int) (domain.Page[domain.Book], error) {
	return u.books.ListCatalog(ctx, cursor, normalizedLimit(limit))
}

func normalizedLimit(limit int) int {
	if limit < 1 {
		return 20
	}
	if limit > 100 {
		return 100
	}
	return limit
}
