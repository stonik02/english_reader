package refresh

import (
	"context"
	"time"

	domain "github.com/deniskrylov/english-reader/backend/internal/domain/auth"
)

type Request struct {
	RefreshToken string
	DeviceLabel  string
}

type UseCase struct {
	find    SessionUserFinder
	revoke  SessionRevoker
	create  SessionCreator
	tokens  TokenIssuer
	refresh RefreshTokenService
	ttl     time.Duration
}

func New(f SessionUserFinder, r SessionRevoker, c SessionCreator, t TokenIssuer, rt RefreshTokenService, ttl time.Duration) *UseCase {
	return &UseCase{
		find:    f,
		revoke:  r,
		create:  c,
		tokens:  t,
		refresh: rt,
		ttl:     ttl,
	}
}
func (u *UseCase) Execute(c context.Context, q Request) (domain.Tokens, error) {
	old := u.refresh.Hash(q.RefreshToken)

	user, e := u.find.FindUserByRefreshHash(c, old)
	if e != nil {
		return domain.Tokens{}, domain.ErrInvalidCredentials
	}

	if e = u.revoke.RevokeSession(c, old); e != nil {
		return domain.Tokens{}, e
	}

	next, h, e := u.refresh.New()
	if e != nil {
		return domain.Tokens{}, e
	}

	if e = u.create.CreateSession(c, user.ID, h, time.Now().Add(u.ttl), q.DeviceLabel); e != nil {
		return domain.Tokens{}, e
	}

	at, exp, e := u.tokens.Issue(user.ID)
	if e != nil {
		return domain.Tokens{}, e
	}

	return domain.Tokens{AccessToken: at, RefreshToken: next, AccessExpiresAt: exp, User: user}, nil
}
