package getme

import (
	"context"

	domain "github.com/deniskrylov/english-reader/backend/internal/domain/auth"
)

type UseCase interface {
	Execute(context.Context, string) (domain.User, error)
}
