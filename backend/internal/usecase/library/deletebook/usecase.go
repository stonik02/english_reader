package deletebook

import (
	"context"
	"errors"
	"fmt"
	"os"

	domain "github.com/deniskrylov/english-reader/backend/internal/domain/library"
)

type UseCase struct {
	books   Books
	storage Storage
}

func New(books Books, storage Storage) *UseCase {
	return &UseCase{books: books, storage: storage}
}

func (u *UseCase) Execute(ctx context.Context, userID, bookID string) error {
	allowed, err := u.books.CanDelete(ctx, userID, bookID)
	if err != nil {
		return err
	}
	if !allowed {
		return domain.ErrForbidden
	}
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
