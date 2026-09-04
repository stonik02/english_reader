package addtomylibrary

import (
	"context"
	"testing"

	domain "github.com/deniskrylov/english-reader/backend/internal/domain/library"
	"go.uber.org/mock/gomock"
)

func TestUseCaseExecuteAddsReadyBook(t *testing.T) {
	controller := gomock.NewController(t)
	mock := NewMockBooks(controller)
	mock.EXPECT().Get(gomock.Any(), "book-1").Return(domain.Book{Status: "ready"}, nil)
	mock.EXPECT().Add(gomock.Any(), "user-1", "book-1", "button").Return(domain.UserBook{}, nil)

	if _, err := New(mock).Execute(context.Background(), "user-1", "book-1"); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
}

func TestUseCaseExecuteRejectsUnreadyBook(t *testing.T) {
	controller := gomock.NewController(t)
	mock := NewMockBooks(controller)
	mock.EXPECT().Get(gomock.Any(), "book-1").Return(domain.Book{Status: "processing"}, nil)

	_, err := New(mock).Execute(context.Background(), "user-1", "book-1")
	if err != domain.ErrNotReady {
		t.Fatalf("Execute() error = %v, want ErrNotReady", err)
	}
}
