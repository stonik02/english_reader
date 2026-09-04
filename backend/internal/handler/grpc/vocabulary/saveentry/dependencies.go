package saveentry

import (
	"context"

	domain "github.com/deniskrylov/english-reader/backend/internal/domain/vocabulary"
	uc "github.com/deniskrylov/english-reader/backend/internal/usecase/vocabulary/saveentry"
)

type UseCase interface {
	Execute(context.Context, string, uc.Request) (domain.Entry, bool, error)
}

type TokenParser interface {
	Parse(string) (string, error)
}
