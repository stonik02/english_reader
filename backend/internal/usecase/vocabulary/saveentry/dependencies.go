package saveentry

import (
	"context"

	domain "github.com/deniskrylov/english-reader/backend/internal/domain/vocabulary"
)

type Repository interface {
	Save(context.Context, string, domain.SaveRequest) (domain.Entry, bool, error)
}
