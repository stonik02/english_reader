package translatetext

import (
	"context"

	domain "github.com/deniskrylov/english-reader/backend/internal/domain/dictionary"
	uc "github.com/deniskrylov/english-reader/backend/internal/usecase/dictionary/translatetext"
)

type UseCase interface {
	Execute(context.Context, uc.Request) (domain.TextTranslationResponse, error)
}

type TokenParser interface{ Parse(string) (string, error) }
