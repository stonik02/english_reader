package getchapter

import (
	"net/http"

	"github.com/deniskrylov/english-reader/backend/internal/handler/http/identity"
	"github.com/deniskrylov/english-reader/backend/internal/handler/http/response"
	"github.com/go-chi/chi/v5"
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
	chapter, err := h.usecase.Execute(r.Context(), userID, chi.URLParam(r, "bookID"), r.URL.Query().Get("chapter_id"), r.URL.Query().Get("href"))
	if err != nil {
		response.ReaderError(w, err)
		return
	}
	response.Write(w, http.StatusOK, chapter)
}
