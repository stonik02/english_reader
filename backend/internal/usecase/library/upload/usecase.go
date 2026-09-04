package upload

import (
	"context"
	"errors"
	"io"
	"path/filepath"

	domain "github.com/deniskrylov/english-reader/backend/internal/domain/library"
	"github.com/google/uuid"
)

type UseCase struct {
	storage Storage
	books   Books
}
type Request struct {
	UserID   string
	Filename string
	File     io.Reader
}

func New(storage Storage, books Books) *UseCase {
	return &UseCase{
		storage: storage,
		books:   books,
	}
}

func (u *UseCase) Execute(ctx context.Context, request Request) (domain.Book, error) {
	stored, err := u.storage.StoreTemporary(request.Filename, request.File)
	if err != nil {
		return domain.Book{}, err
	}

	defer u.storage.Remove(stored.Temp)

	if book, err := u.books.FindBySHA256(ctx, stored.SHA256[:]); err == nil {
		return book, nil
	} else if !errors.Is(err, domain.ErrNotFound) {
		return domain.Book{}, err
	}

	id := uuid.NewString()

	path, err := u.storage.MoveToBook(stored.Temp, id)
	if err != nil {
		return domain.Book{}, err
	}

	book, err := u.books.Create(ctx, id, request.UserID, filepath.Base(request.Filename), path, stored.SHA256[:])
	if err != nil {
		u.storage.Remove(path)
		return domain.Book{}, err
	}

	return book, nil
}
