package addtomylibrary

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

func (u *UseCase) Execute(ctx context.Context, userID, bookID string) (domain.UserBook, error) {
	book, err := u.books.Get(ctx, bookID)
	if err != nil {
		return domain.UserBook{}, err
	}
	if book.Status != "ready" {
		return domain.UserBook{}, domain.ErrNotReady
	}
	return u.books.Add(ctx, userID, bookID, "button")
}
