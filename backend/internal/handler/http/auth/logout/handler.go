package logout

import "net/http"

type Handler struct {
	usecase UseCase
	secure  bool
}

func New(u UseCase, s bool) *Handler { return &Handler{u, s} }
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if c, e := r.Cookie("refresh_token"); e == nil {
		_ = h.usecase.Execute(r.Context(), c.Value)
	}
	http.SetCookie(w, &http.Cookie{Name: "refresh_token", Value: "", Path: "/api/v1/auth", MaxAge: -1, HttpOnly: true, Secure: h.secure, SameSite: http.SameSiteStrictMode})
	w.WriteHeader(204)
}
