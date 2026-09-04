package deleteentry

import (
	"context"

	uc "github.com/deniskrylov/english-reader/backend/internal/usecase/vocabulary/deleteentry"
)

type UseCase interface {
	Execute(context.Context, string, uc.Request) error
}

type TokenParser interface {
	Parse(string) (string, error)
}
