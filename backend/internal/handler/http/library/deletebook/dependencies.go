package deletebook

import "context"

type UseCase interface {
	Execute(context.Context, string) error
}
type TokenParser interface{ Parse(string) (string, error) }
