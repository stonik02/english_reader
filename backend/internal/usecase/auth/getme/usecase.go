package getme

import (
	"context"

	domain "github.com/deniskrylov/english-reader/backend/internal/domain/auth"
)

type UseCase struct {
	tokens AccessTokenParser
	users  UserFinder
}

func New(t AccessTokenParser, u UserFinder) *UseCase {
	return &UseCase{
		t, u,
	}
}

func (u *UseCase) Execute(c context.Context, access string) (domain.User, error) {
	id, e := u.tokens.Parse(access)
	if e != nil {
		return domain.User{}, domain.ErrInvalidCredentials
	}

	user, e := u.users.FindUserByID(c, id)
	if e != nil {
		return domain.User{}, domain.ErrInvalidCredentials
	}

	return user, nil
}
