package listmybooks

import (
	"net/http"
	"strconv"

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
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	page, err := h.usecase.Execute(r.Context(), userID, r.URL.Query().Get("cursor"), limit)
	if err != nil {
		response.LibraryError(w, err)
		return
	}
	response.Write(w, http.StatusOK, page)
}
