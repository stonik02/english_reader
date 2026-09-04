package logout

import "context"

type UseCase struct {
	sessions SessionRevoker
	refresh  RefreshTokenHasher
}

func New(s SessionRevoker, r RefreshTokenHasher) *UseCase {
	return &UseCase{
		sessions: s,
		refresh:  r,
	}
}

func (u *UseCase) Execute(c context.Context, token string) error {
	return u.sessions.RevokeSession(c, u.refresh.Hash(token))
}
