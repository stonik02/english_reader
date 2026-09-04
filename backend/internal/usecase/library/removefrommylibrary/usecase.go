package removefrommylibrary

import "context"

type UseCase struct {
	books Books
}

func New(books Books) *UseCase {
	return &UseCase{books: books}
}

func (u *UseCase) Execute(ctx context.Context, userID, bookID string) error {
	return u.books.Remove(ctx, userID, bookID)
}
