package uploadbook

import (
	"context"

	domain "github.com/deniskrylov/english-reader/backend/internal/domain/library"
	uc "github.com/deniskrylov/english-reader/backend/internal/usecase/library/upload"
)

type UseCase interface {
	Execute(context.Context, uc.Request) (domain.Book, error)
}

type TokenParser interface {
	Parse(string) (string, error)
}
