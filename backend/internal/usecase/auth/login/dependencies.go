package login

import (
	"context"
	"time"

	domain "github.com/deniskrylov/english-reader/backend/internal/domain/auth"
)

type AccountFinder interface {
	FindUserByEmail(context.Context, string) (domain.User, string, error)
}

type SessionCreator interface {
	CreateSession(context.Context, string, []byte, time.Time, string) error
}

type PasswordVerifier interface {
	Verify(string, string) bool
}

type TokenIssuer interface {
	Issue(string) (string, time.Time, error)
}

type RefreshTokenGenerator interface {
	New() (string, []byte, error)
}
