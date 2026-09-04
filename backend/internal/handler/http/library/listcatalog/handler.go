package listcatalog

import (
	"net/http"
	"strconv"

	"github.com/deniskrylov/english-reader/backend/internal/handler/http/response"
)

type Handler struct {
	usecase UseCase
}

func New(usecase UseCase) *Handler {
	return &Handler{usecase: usecase}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	page, err := h.usecase.Execute(r.Context(), r.URL.Query().Get("cursor"), limit)
	if err != nil {
		response.LibraryError(w, err)
		return
	}
	response.Write(w, http.StatusOK, page)
}
