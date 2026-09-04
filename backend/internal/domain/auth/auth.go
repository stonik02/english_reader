package auth

import (
	"errors"
	"time"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailTaken         = errors.New("email is already registered")
	ErrInvalidInput       = errors.New("invalid input")
)

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}
type Tokens struct {
	AccessToken     string
	RefreshToken    string
	AccessExpiresAt time.Time
	User            User
}
