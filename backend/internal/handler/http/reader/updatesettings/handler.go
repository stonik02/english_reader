package updatesettings

import (
	"encoding/json"
	domain "github.com/deniskrylov/english-reader/backend/internal/domain/reader"
	"github.com/deniskrylov/english-reader/backend/internal/handler/http/identity"
	"github.com/deniskrylov/english-reader/backend/internal/handler/http/response"
	"net/http"
)

type Handler struct {
	usecase UseCase
	tokens  TokenParser
}

func New(usecase UseCase, tokens TokenParser) *Handler {
	return &Handler{usecase: usecase, tokens: tokens}
}
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userID, err := identity.Subject(r, h.tokens)
	if err != nil {
		response.Write(w, http.StatusUnauthorized, map[string]string{"error": "invalid access token"})
		return
	}
	var value domain.Settings
	if err := json.NewDecoder(r.Body).Decode(&value); err != nil {
		response.Write(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}
	value, err = h.usecase.Execute(r.Context(), userID, value)
	if err != nil {
		response.ReaderError(w, err)
		return
	}
	response.Write(w, http.StatusOK, value)
}
