package logout

import "context"

type UseCase interface {
	Execute(context.Context, string) error
}
