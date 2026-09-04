package response

import (
	"errors"
	"time"

	readerv1 "github.com/deniskrylov/english-reader/backend/gen/reader/v1"
	domain "github.com/deniskrylov/english-reader/backend/internal/domain/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Auth(tokens domain.Tokens) *readerv1.AuthResponse {
	return &readerv1.AuthResponse{
		User:                   User(tokens.User),
		AccessToken:            tokens.AccessToken,
		RefreshToken:           tokens.RefreshToken,
		AccessTokenExpiresUnix: tokens.AccessExpiresAt.Unix(),
	}
}

func User(user domain.User) *readerv1.User {
	return &readerv1.User{
		Id:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.UTC().Format(time.RFC3339),
	}
}

func Error(err error) error {
	switch {
	case errors.Is(err, domain.ErrInvalidInput):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, domain.ErrEmailTaken):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, domain.ErrInvalidCredentials):
		return status.Error(codes.Unauthenticated, err.Error())
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
