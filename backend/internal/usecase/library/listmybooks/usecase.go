package listmybooks

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

func (u *UseCase) Execute(ctx context.Context, userID, cursor string, limit int) (domain.Page[domain.UserBook], error) {
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	return u.books.ListMine(ctx, userID, cursor, limit)
}
