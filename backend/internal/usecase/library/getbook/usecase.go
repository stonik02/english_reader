package getbook

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

func (u *UseCase) Execute(ctx context.Context, bookID string) (domain.Book, error) {
	return u.books.Get(ctx, bookID)
}
