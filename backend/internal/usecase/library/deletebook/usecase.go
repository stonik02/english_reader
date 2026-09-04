package deletebook

import (
	"context"
	"errors"
	"fmt"
	"os"
)

type UseCase struct {
	books   Books
	storage Storage
}

func New(books Books, storage Storage) *UseCase {
	return &UseCase{books: books, storage: storage}
}

func (u *UseCase) Execute(ctx context.Context, bookID string) error {
	files, err := u.books.Delete(ctx, bookID)
	if err != nil {
		return err
	}
	for _, path := range []string{files.SourcePath, files.CoverPath} {
		if path == "" {
			continue
		}
		if err := u.storage.Delete(path); err != nil && !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("remove book storage: %w", err)
		}
	}
	return nil
}
