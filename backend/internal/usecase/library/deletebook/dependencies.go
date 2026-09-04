package deletebook

import (
	"context"

	domain "github.com/deniskrylov/english-reader/backend/internal/domain/library"
)

type Books interface {
	Delete(context.Context, string) (domain.StoredBookFiles, error)
}

type Storage interface {
	Delete(string) error
}
