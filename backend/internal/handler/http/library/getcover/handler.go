package getcover

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
	if _, err := identity.Subject(r, h.tokens); err != nil {
		response.Write(w, http.StatusUnauthorized, map[string]string{"error": "invalid access token"})
		return
	}
	data, contentType, err := h.usecase.Execute(r.Context(), chi.URLParam(r, "bookID"))
	if err != nil {
		response.LibraryError(w, err)
		return
	}
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Cache-Control", "private, max-age=3600")
	_, _ = w.Write(data)
}
