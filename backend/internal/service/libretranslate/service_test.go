package libretranslate

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestServiceTranslate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"translatedText":"Привет"}`))
	}))
	defer server.Close()
	value, err := New(server.URL, time.Second, 100).Translate(context.Background(), "Hello")
	if err != nil || value != "Привет" {
		t.Fatalf("Translate()=%q,%v", value, err)
	}
}
