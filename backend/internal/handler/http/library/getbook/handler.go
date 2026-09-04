package getbook

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/deniskrylov/english-reader/backend/internal/handler/http/response"
)

type Handler struct {
	usecase UseCase
}

func New(usecase UseCase) *Handler {
	return &Handler{usecase: usecase}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	book, err := h.usecase.Execute(r.Context(), chi.URLParam(r, "bookID"))
	if err != nil {
		response.LibraryError(w, err)
		return
	}
	response.Write(w, http.StatusOK, book)
}
