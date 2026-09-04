package response

import (
	"encoding/json"
	"errors"
	"net/http"

	domain "github.com/deniskrylov/english-reader/backend/internal/domain/auth"
	library "github.com/deniskrylov/english-reader/backend/internal/domain/library"
	reader "github.com/deniskrylov/english-reader/backend/internal/domain/reader"
)

func Write(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func ReaderError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, reader.ErrInvalidInput):
		Write(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
	case errors.Is(err, reader.ErrNotFound):
		Write(w, http.StatusNotFound, map[string]string{"error": err.Error()})
	case errors.Is(err, reader.ErrNotReady):
		Write(w, http.StatusPreconditionFailed, map[string]string{"error": err.Error()})
	default:
		Write(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}
}

func LibraryError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, library.ErrNotFound):
		Write(w, http.StatusNotFound, map[string]string{"error": err.Error()})
	case errors.Is(err, library.ErrNotReady):
		Write(w, http.StatusPreconditionFailed, map[string]string{"error": err.Error()})
	case errors.Is(err, library.ErrInvalidUpload):
		Write(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
	case errors.Is(err, library.ErrTooLarge):
		Write(w, http.StatusRequestEntityTooLarge, map[string]string{"error": err.Error()})
	default:
		Write(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}
}

func Error(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidInput):
		Write(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
	case errors.Is(err, domain.ErrEmailTaken):
		Write(w, http.StatusConflict, map[string]string{"error": err.Error()})
	case errors.Is(err, domain.ErrInvalidCredentials):
		Write(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
	default:
		Write(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}
}

func InvalidJSON(w http.ResponseWriter) {
	Write(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
}
