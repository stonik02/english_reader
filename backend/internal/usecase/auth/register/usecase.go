package register

import (
	"context"
	"net/mail"
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
	accounts  AccountCreator
	passwords PasswordHasher
	tokens    TokenIssuer
	refresh   RefreshTokenGenerator
	ttl       time.Duration
}

func New(a AccountCreator, p PasswordHasher, t TokenIssuer, r RefreshTokenGenerator, ttl time.Duration) *UseCase {
	return &UseCase{
		accounts:  a,
		passwords: p,
		tokens:    t,
		refresh:   r,
		ttl:       ttl,
	}
}

func (u *UseCase) Execute(c context.Context, q Request) (domain.Tokens, error) {
	email := strings.ToLower(strings.TrimSpace(q.Email))

	parsed, e := mail.ParseAddress(email)
	if e != nil || parsed.Address != email || len(q.Password) < 12 || len(q.Password) > 256 {
		return domain.Tokens{}, domain.ErrInvalidInput
	}

	pass, e := u.passwords.Hash(q.Password)
	if e != nil {
		return domain.Tokens{}, e
	}

	rt, h, e := u.refresh.New()
	if e != nil {
		return domain.Tokens{}, e
	}

	user, e := u.accounts.CreateUserWithSession(c, email, pass, h, time.Now().Add(u.ttl), q.DeviceLabel)
	if e != nil {
		return domain.Tokens{}, e
	}

	at, exp, e := u.tokens.Issue(user.ID)
	if e != nil {
		return domain.Tokens{}, e
	}

	return domain.Tokens{AccessToken: at, RefreshToken: rt, AccessExpiresAt: exp, User: user}, nil
}
