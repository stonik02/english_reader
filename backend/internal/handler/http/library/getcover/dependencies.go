package getcover

import "context"

type UseCase interface {
	Execute(context.Context, string) ([]byte, string, error)
}

type TokenParser interface {
	Parse(string) (string, error)
}
