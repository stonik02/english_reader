package lookupword

import (
	"context"
	domain "github.com/deniskrylov/english-reader/backend/internal/domain/dictionary"
	uc "github.com/deniskrylov/english-reader/backend/internal/usecase/dictionary/lookupword"
)

type UseCase interface {
	Execute(context.Context, uc.Request) (domain.LookupResponse, error)
}
type TokenParser interface{ Parse(string) (string, error) }
