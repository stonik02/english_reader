package grantadmin

import (
	"context"

	domain "github.com/deniskrylov/english-reader/backend/internal/domain/auth"
)

type UserFinder interface {
	FindUserByEmail(context.Context, string) (domain.User, string, error)
}

type RoleUpdater interface {
	SetRole(context.Context, string, string) error
}
