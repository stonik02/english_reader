package logout

import "context"

type SessionRevoker interface {
	RevokeSession(context.Context, []byte) error
}

type RefreshTokenHasher interface {
	Hash(string) []byte
}
