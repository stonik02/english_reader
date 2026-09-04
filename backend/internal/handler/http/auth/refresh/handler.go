package refresh

import (
	"net/http"
	"time"

	"github.com/deniskrylov/english-reader/backend/internal/handler/http/response"
	uc "github.com/deniskrylov/english-reader/backend/internal/usecase/auth/refresh"
)

type Handler struct {
	usecase UseCase
	secure  bool
	ttl     time.Duration
}

func New(u UseCase, s bool, t time.Duration) *Handler { return &Handler{u, s, t} }
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, e := r.Cookie("refresh_token")
	if e != nil {
		response.Write(w, http.StatusUnauthorized, map[string]string{"error": "invalid refresh token"})
		return
	}
	v, e := h.usecase.Execute(r.Context(), uc.Request{RefreshToken: c.Value, DeviceLabel: r.UserAgent()})
	if e != nil {
		response.Error(w, e)
		return
	}
	http.SetCookie(w, &http.Cookie{Name: "refresh_token", Value: v.RefreshToken, Path: "/api/v1/auth", HttpOnly: true, Secure: h.secure, SameSite: http.SameSiteStrictMode, Expires: time.Now().Add(h.ttl)})
	response.Write(w, http.StatusOK, map[string]any{"user": v.User, "access_token": v.AccessToken, "access_token_expires_at": v.AccessExpiresAt})
}
