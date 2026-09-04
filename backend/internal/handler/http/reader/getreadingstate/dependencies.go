package getreadingstate

import (
	"context"

	domain "github.com/deniskrylov/english-reader/backend/internal/domain/reader"
)

type UseCase interface {
	Execute(context.Context, string, string) (domain.State, error)
}

type TokenParser interface {
	Parse(string) (string, error)
}
