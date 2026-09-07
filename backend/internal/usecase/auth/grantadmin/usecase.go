package grantadmin

import (
	"context"
	"strings"

	domain "github.com/deniskrylov/english-reader/backend/internal/domain/auth"
)

type UseCase struct {
	users UserFinder
	roles RoleUpdater
}

func New(users UserFinder, roles RoleUpdater) *UseCase {
	return &UseCase{users: users, roles: roles}
}

func (u *UseCase) Execute(ctx context.Context, email string) (domain.User, error) {
	user, _, err := u.users.FindUserByEmail(ctx, strings.ToLower(strings.TrimSpace(email)))
	if err != nil {
		return domain.User{}, err
	}
	if err := u.roles.SetRole(ctx, user.ID, "admin"); err != nil {
		return domain.User{}, err
	}
	user.Role = "admin"
	return user, nil
}
