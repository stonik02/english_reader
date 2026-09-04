package getcover

import (
	"context"
	"fmt"
	"path/filepath"

	domain "github.com/deniskrylov/english-reader/backend/internal/domain/library"
)

type UseCase struct {
	books  Books
	covers Covers
}

func New(books Books, covers Covers) *UseCase {
	return &UseCase{books: books, covers: covers}
}

func (u *UseCase) Execute(ctx context.Context, bookID string) ([]byte, string, error) {
	book, err := u.books.Get(ctx, bookID)
	if err != nil {
		return nil, "", err
	}
	if book.CoverPath == "" {
		return nil, "", domain.ErrNotFound
	}
	contentType := map[string]string{".jpg": "image/jpeg", ".png": "image/png", ".gif": "image/gif", ".webp": "image/webp"}[filepath.Ext(book.CoverPath)]
	if contentType == "" {
		return nil, "", fmt.Errorf("unsupported stored cover type")
	}
	data, err := u.covers.Read(book.CoverPath)
	return data, contentType, err
}
