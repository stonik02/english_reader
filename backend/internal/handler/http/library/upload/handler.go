package upload

import (
	"net/http"

	"github.com/deniskrylov/english-reader/backend/internal/handler/http/identity"
	"github.com/deniskrylov/english-reader/backend/internal/handler/http/response"
	uc "github.com/deniskrylov/english-reader/backend/internal/usecase/library/upload"
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
	if err := r.ParseMultipartForm(1 << 20); err != nil {
		response.Write(w, http.StatusBadRequest, map[string]string{"error": "invalid multipart form"})
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		response.Write(w, http.StatusBadRequest, map[string]string{"error": "EPUB file is required"})
		return
	}
	defer file.Close()
	book, err := h.usecase.Execute(r.Context(), uc.Request{UserID: userID, Filename: header.Filename, File: file})
	if err != nil {
		response.LibraryError(w, err)
		return
	}
	response.Write(w, http.StatusCreated, book)
}
