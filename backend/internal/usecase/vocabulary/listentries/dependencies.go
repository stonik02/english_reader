package listentries

import (
	"context"

	domain "github.com/deniskrylov/english-reader/backend/internal/domain/vocabulary"
)

type Repository interface {
	List(context.Context, string, domain.ListRequest) ([]domain.Entry, error)
}

type WordNormalizer interface {
	Normalize(string) (string, error)
}
