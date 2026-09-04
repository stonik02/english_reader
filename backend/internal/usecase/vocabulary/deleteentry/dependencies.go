package deleteentry

import "context"

type Repository interface {
	DeleteByID(context.Context, string, string) error
	DeleteByLemmaID(context.Context, string, int64) error
}
