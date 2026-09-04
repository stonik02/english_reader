package identity

type TokenParser interface{ Parse(string) (string, error) }
