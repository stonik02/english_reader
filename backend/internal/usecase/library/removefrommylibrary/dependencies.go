package removefrommylibrary

import "context"

type Books interface {
	Remove(context.Context, string, string) error
}
