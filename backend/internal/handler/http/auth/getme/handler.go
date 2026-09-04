package getme

import (
	"net/http"
	"strings"

	"github.com/deniskrylov/english-reader/backend/internal/handler/http/response"
)

type Handler struct{ usecase UseCase }

func New(u UseCase) *Handler { return &Handler{u} }
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v, e := h.usecase.Execute(r.Context(), strings.TrimSpace(strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")))
	if e != nil {
		response.Error(w, e)
		return
	}
	response.Write(w, http.StatusOK, v)
}
