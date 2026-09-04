package getsettings

import (
	"context"

	domain "github.com/deniskrylov/english-reader/backend/internal/domain/reader"
)

type UseCase interface {
	Execute(context.Context, string) (domain.Settings, error)
}
type TokenParser interface{ Parse(string) (string, error) }
