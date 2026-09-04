package token

import (
	"time"

	"github.com/deniskrylov/english-reader/backend/internal/domain/auth"
	"github.com/golang-jwt/jwt/v5"
)

type Service struct {
	secret []byte
	ttl    time.Duration
}

func New(secret string, ttl time.Duration) *Service { return &Service{[]byte(secret), ttl} }

func (s *Service) Issue(subject string) (string, time.Time, error) {
	now := time.Now()
	exp := now.Add(s.ttl)

	claims := jwt.RegisteredClaims{Subject: subject, ExpiresAt: jwt.NewNumericDate(exp), IssuedAt: jwt.NewNumericDate(now), Issuer: "english-reader"}

	v, e := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.secret)

	return v, exp, e
}
func (s *Service) Parse(v string) (string, error) {
	c := &jwt.RegisteredClaims{}

	t, e := jwt.ParseWithClaims(v, c, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, auth.ErrInvalidCredentials
		}
		return s.secret, nil
	})

	if e != nil || !t.Valid || c.Subject == "" {
		return "", auth.ErrInvalidCredentials
	}

	return c.Subject, nil
}
