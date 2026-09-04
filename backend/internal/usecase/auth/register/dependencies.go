package register

import (
	"context"
	"time"

	domain "github.com/deniskrylov/english-reader/backend/internal/domain/auth"
)

type AccountCreator interface {
	CreateUserWithSession(context.Context, string, string, []byte, time.Time, string) (domain.User, error)
}

type PasswordHasher interface {
	Hash(string) (string, error)
}

type TokenIssuer interface {
	Issue(string) (string, time.Time, error)
}

type RefreshTokenGenerator interface {
	New() (string, []byte, error)
}
