package getreadingstate

import (
	"net/http"
	"net/http/httptest"
	"testing"

	domain "github.com/deniskrylov/english-reader/backend/internal/domain/reader"
	"github.com/go-chi/chi/v5"
	"go.uber.org/mock/gomock"
)

func TestHandlerServeHTTP(t *testing.T) {
	controller := gomock.NewController(t)
	usecase := NewMockUseCase(controller)
	tokens := NewMockTokenParser(controller)
	tokens.EXPECT().Parse("access-token").Return("user-1", nil)
	usecase.EXPECT().Execute(gomock.Any(), "user-1", "book-1").Return(domain.State{}, nil)

	router := chi.NewRouter()
	router.Get("/books/{bookID}/state", New(usecase, tokens).ServeHTTP)
	request := httptest.NewRequest(http.MethodGet, "/books/book-1/state", nil)
	request.Header.Set("Authorization", "Bearer access-token")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}
}
