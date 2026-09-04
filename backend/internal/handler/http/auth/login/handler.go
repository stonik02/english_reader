package login

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/deniskrylov/english-reader/backend/internal/handler/http/response"
	uc "github.com/deniskrylov/english-reader/backend/internal/usecase/auth/login"
)

type Handler struct {
	usecase UseCase
	secure  bool
	ttl     time.Duration
}

func New(u UseCase, s bool, t time.Duration) *Handler { return &Handler{u, s, t} }
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var q uc.Request
	if e := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20)).Decode(&q); e != nil {
		response.InvalidJSON(w)
		return
	}
	v, e := h.usecase.Execute(r.Context(), q)
	if e != nil {
		response.Error(w, e)
		return
	}
	http.SetCookie(w, &http.Cookie{Name: "refresh_token", Value: v.RefreshToken, Path: "/api/v1/auth", HttpOnly: true, Secure: h.secure, SameSite: http.SameSiteStrictMode, Expires: time.Now().Add(h.ttl)})
	response.Write(w, http.StatusOK, map[string]any{"user": v.User, "access_token": v.AccessToken, "access_token_expires_at": v.AccessExpiresAt})
}
