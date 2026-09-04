package identity

import (
	"net/http"
	"strings"
)

func Subject(request *http.Request, parser TokenParser) (string, error) {
	token := strings.TrimSpace(strings.TrimPrefix(request.Header.Get("Authorization"), "Bearer "))
	return parser.Parse(token)
}
