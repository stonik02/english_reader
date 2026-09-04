package saveprogress

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/deniskrylov/english-reader/backend/internal/handler/http/identity"
	"github.com/deniskrylov/english-reader/backend/internal/handler/http/response"
	uc "github.com/deniskrylov/english-reader/backend/internal/usecase/reader/saveprogress"
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

	var request uc.Request
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response.Write(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}

	progress, err := h.usecase.Execute(r.Context(), userID, chi.URLParam(r, "bookID"), request)
	if err != nil {
		response.ReaderError(w, err)
		return
	}

	response.Write(w, http.StatusOK, progress)
}
