package getchapter

import (
	"context"

	domain "github.com/deniskrylov/english-reader/backend/internal/domain/reader"
)

type UseCase interface {
	Execute(context.Context, string, string, string, string) (domain.Chapter, error)
}
type TokenParser interface{ Parse(string) (string, error) }
