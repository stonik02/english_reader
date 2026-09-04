package removefrommylibrary

import "context"

type UseCase interface {
	Execute(context.Context, string, string) error
}

type TokenParser interface {
	Parse(string) (string, error)
}
