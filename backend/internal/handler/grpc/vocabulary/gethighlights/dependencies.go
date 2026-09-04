package gethighlights

import (
	"context"

	domain "github.com/deniskrylov/english-reader/backend/internal/domain/vocabulary"
	uc "github.com/deniskrylov/english-reader/backend/internal/usecase/vocabulary/gethighlights"
)

type UseCase interface {
	Execute(context.Context, uc.Request) ([]domain.HighlightToken, error)
}

type TokenParser interface {
	Parse(string) (string, error)
}
