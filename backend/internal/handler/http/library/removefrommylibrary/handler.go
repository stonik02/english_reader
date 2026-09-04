package removefrommylibrary

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/deniskrylov/english-reader/backend/internal/handler/http/identity"
	"github.com/deniskrylov/english-reader/backend/internal/handler/http/response"
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
	if err := h.usecase.Execute(r.Context(), userID, chi.URLParam(r, "bookID")); err != nil {
		response.LibraryError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
