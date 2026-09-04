package listmybooks

import (
	"context"

	domain "github.com/deniskrylov/english-reader/backend/internal/domain/library"
)

type UseCase interface {
	Execute(context.Context, string, string, int) (domain.Page[domain.UserBook], error)
}

type TokenParser interface {
	Parse(string) (string, error)
}
