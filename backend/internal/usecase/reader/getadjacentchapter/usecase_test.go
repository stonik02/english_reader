package getadjacentchapter

import (
	"context"
	"testing"

	domain "github.com/deniskrylov/english-reader/backend/internal/domain/reader"
	"go.uber.org/mock/gomock"
)

func TestUseCaseRejectsInvalidDirection(t *testing.T) {
	controller := gomock.NewController(t)
	_, err := New(NewMockReaderRepository(controller)).Execute(context.Background(), "user-1", "book-1", "chapter-1", 0)
	if err != domain.ErrInvalidInput {
		t.Fatalf("error = %v", err)
	}
}

func TestUseCaseExecute(t *testing.T) {
	controller := gomock.NewController(t)
	repository := NewMockReaderRepository(controller)
	repository.EXPECT().EnsureReadyBook(gomock.Any(), "user-1", "book-1").Return(nil)
	repository.EXPECT().Adjacent(gomock.Any(), "book-1", "chapter-1", 1).Return(domain.Chapter{ID: "chapter-2"}, nil)
	if _, err := New(repository).Execute(context.Background(), "user-1", "book-1", "chapter-1", 1); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
}
