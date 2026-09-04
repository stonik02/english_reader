package getme

import (
	"context"

	domain "github.com/deniskrylov/english-reader/backend/internal/domain/auth"
)

type AccessTokenParser interface {
	Parse(string) (string, error)
}

type UserFinder interface {
	FindUserByID(context.Context, string) (domain.User, error)
}
