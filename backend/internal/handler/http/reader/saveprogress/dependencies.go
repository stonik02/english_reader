package saveprogress

import (
	"context"

	domain "github.com/deniskrylov/english-reader/backend/internal/domain/reader"
	uc "github.com/deniskrylov/english-reader/backend/internal/usecase/reader/saveprogress"
)

type UseCase interface {
	Execute(context.Context, string, string, uc.Request) (domain.Progress, error)
}

type TokenParser interface {
	Parse(string) (string, error)
}
