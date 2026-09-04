package listentries

import (
	"context"

	uc "github.com/deniskrylov/english-reader/backend/internal/usecase/vocabulary/listentries"
)

type UseCase interface {
	Execute(context.Context, string, uc.Request) (uc.Response, error)
}

type TokenParser interface {
	Parse(string) (string, error)
}
