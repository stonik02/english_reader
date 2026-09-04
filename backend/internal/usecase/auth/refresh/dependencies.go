package refresh

import (
	"context"
	"time"

	domain "github.com/deniskrylov/english-reader/backend/internal/domain/auth"
)

type SessionUserFinder interface {
	FindUserByRefreshHash(context.Context, []byte) (domain.User, error)
}

type SessionRevoker interface {
	RevokeSession(context.Context, []byte) error
}

type SessionCreator interface {
	CreateSession(context.Context, string, []byte, time.Time, string) error
}

type TokenIssuer interface {
	Issue(string) (string, time.Time, error)
}

type RefreshTokenService interface {
	New() (string, []byte, error)
	Hash(string) []byte
}
