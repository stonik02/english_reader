package register

import (
	"context"
	domain "github.com/deniskrylov/english-reader/backend/internal/domain/auth"
	uc "github.com/deniskrylov/english-reader/backend/internal/usecase/auth/register"
)

type UseCase interface {
	Execute(context.Context, uc.Request) (domain.Tokens, error)
}
