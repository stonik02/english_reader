package listmybooks

import (
	"context"

	domain "github.com/deniskrylov/english-reader/backend/internal/domain/library"
)

type Books interface {
	ListMine(context.Context, string, string, int) (domain.Page[domain.UserBook], error)
}
