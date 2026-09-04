package getreadingstate

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

func (h *Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	userID, err := identity.Subject(request, h.tokens)
	if err != nil {
		response.Write(writer, http.StatusUnauthorized, map[string]string{"error": "invalid access token"})
		return
	}

	state, err := h.usecase.Execute(request.Context(), userID, chi.URLParam(request, "bookID"))
	if err != nil {
		response.ReaderError(writer, err)
		return
	}

	response.Write(writer, http.StatusOK, state)
}
