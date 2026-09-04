package httpserver

import (
	"context"
	_ "embed"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	httpSwagger "github.com/swaggo/http-swagger"
)

//go:embed openapi.yaml
var openAPISpec []byte

func New(pool *pgxpool.Pool, frontendOrigin string, register, login, refresh, logout, me, upload, listCatalog, getBook, getCover, deleteBook, addToMyLibrary, listMyBooks, removeFromMyLibrary, getReadingState, getChapter, getAdjacentChapter, saveProgress, getSettings, updateSettings http.Handler) http.Handler {
	r := chi.NewRouter()
	r.Use(recoverer)
	r.Use(cors(frontendOrigin))
	r.Get("/health/live", live)
	r.Get("/health/ready", ready(pool))
	r.Get("/openapi.yaml", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/yaml")
		_, _ = w.Write(openAPISpec)
	})
	r.Mount("/swagger", httpSwagger.Handler(httpSwagger.URL("/openapi.yaml")))
	r.Route("/api/v1/auth", func(a chi.Router) {
		a.Post("/register", register.ServeHTTP)
		a.Post("/login", login.ServeHTTP)
		a.Post("/refresh", refresh.ServeHTTP)
		a.Post("/logout", logout.ServeHTTP)
		a.Get("/me", me.ServeHTTP)
	})
	r.Route("/api/v1/library", func(l chi.Router) {
		l.Post("/books", upload.ServeHTTP)
		l.Get("/books", listCatalog.ServeHTTP)
		l.Get("/books/{bookID}", getBook.ServeHTTP)
		l.Get("/books/{bookID}/cover", getCover.ServeHTTP)
		l.Delete("/books/{bookID}", deleteBook.ServeHTTP)
		l.Post("/books/{bookID}/my-library", addToMyLibrary.ServeHTTP)
		l.Get("/my-books", listMyBooks.ServeHTTP)
		l.Delete("/books/{bookID}/my-library", removeFromMyLibrary.ServeHTTP)
	})
	r.Route("/api/v1/reader", func(reader chi.Router) {
		reader.Get("/books/{bookID}/state", getReadingState.ServeHTTP)
		reader.Get("/books/{bookID}/chapters", getChapter.ServeHTTP)
		reader.Get("/books/{bookID}/chapters/{chapterID}/adjacent", getAdjacentChapter.ServeHTTP)
		reader.Put("/books/{bookID}/progress", saveProgress.ServeHTTP)
		reader.Get("/settings", getSettings.ServeHTTP)
		reader.Put("/settings", updateSettings.ServeHTTP)
	})
	return r
}
func live(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, 200, map[string]string{"status": "ok"})
}

func ready(p *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, cancel := context.WithTimeout(r.Context(), time.Second)
		defer cancel()
		if p.Ping(c) != nil {
			writeJSON(w, 503, map[string]string{"status": "not_ready"})
			return
		}
		writeJSON(w, 200, map[string]string{"status": "ready"})
	}
}

func recoverer(n http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recover() != nil {
				writeJSON(w, 500, map[string]string{"error": "internal server error"})
			}
		}()
		n.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, c int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(c)
	_ = json.NewEncoder(w).Encode(v)
}
