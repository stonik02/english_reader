package login

import (
	"context"
	"strings"
	"time"

	domain "github.com/deniskrylov/english-reader/backend/internal/domain/auth"
)

type Request struct {
	Email       string
	Password    string
	DeviceLabel string
}
type UseCase struct {
	accounts  AccountFinder
	sessions  SessionCreator
	passwords PasswordVerifier
	tokens    TokenIssuer
	refresh   RefreshTokenGenerator
	ttl       time.Duration
}

func New(a AccountFinder, s SessionCreator, p PasswordVerifier, t TokenIssuer, r RefreshTokenGenerator, ttl time.Duration) *UseCase {
	return &UseCase{
		accounts:  a,
		sessions:  s,
		passwords: p,
		tokens:    t,
		refresh:   r,
		ttl:       ttl,
	}
}
func (u *UseCase) Execute(c context.Context, q Request) (domain.Tokens, error) {
	user, hash, e := u.accounts.FindUserByEmail(c, strings.ToLower(strings.TrimSpace(q.Email)))
	if e != nil || !u.passwords.Verify(hash, q.Password) {
		return domain.Tokens{}, domain.ErrInvalidCredentials
	}

	rt, h, e := u.refresh.New()
	if e != nil {
		return domain.Tokens{}, e
	}

	if e = u.sessions.CreateSession(c, user.ID, h, time.Now().Add(u.ttl), q.DeviceLabel); e != nil {
		return domain.Tokens{}, e
	}

	at, exp, e := u.tokens.Issue(user.ID)
	if e != nil {
		return domain.Tokens{}, e
	}

	return domain.Tokens{AccessToken: at, RefreshToken: rt, AccessExpiresAt: exp, User: user}, nil
}
