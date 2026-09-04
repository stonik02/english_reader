package upload

import (
	"context"
	"io"

	domain "github.com/deniskrylov/english-reader/backend/internal/domain/library"
	"github.com/deniskrylov/english-reader/backend/internal/service/epubstorage"
)

type Storage interface {
	StoreTemporary(string, io.Reader) (epubstorage.StoredFile, error)
	MoveToBook(string, string) (string, error)
	Remove(string)
}

type Books interface {
	FindBySHA256(context.Context, []byte) (domain.Book, error)
	Create(context.Context, string, string, string, string, []byte) (domain.Book, error)
	Get(context.Context, string) (domain.Book, error)
}
