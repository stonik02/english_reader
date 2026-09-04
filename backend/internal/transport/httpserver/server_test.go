package httpserver

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLiveEndpoint(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/health/live", nil)
	response := httptest.NewRecorder()

	live(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", response.Code, http.StatusOK)
	}
}
