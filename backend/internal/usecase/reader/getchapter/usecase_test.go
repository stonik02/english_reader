package getchapter

import (
	"context"
	"testing"

	domain "github.com/deniskrylov/english-reader/backend/internal/domain/reader"
	"go.uber.org/mock/gomock"
)

func TestUseCaseExecute(t *testing.T) {
	controller := gomock.NewController(t)
	repository := NewMockReaderRepository(controller)
	repository.EXPECT().EnsureReadyBook(gomock.Any(), "user-1", "book-1").Return(nil)
	repository.EXPECT().Chapter(gomock.Any(), "book-1", "chapter-1", "").Return(domain.Chapter{ID: "chapter-1"}, nil)

	chapter, err := New(repository).Execute(context.Background(), "user-1", "book-1", "chapter-1", "")
	if err != nil || chapter.ID != "chapter-1" {
		t.Fatalf("Execute() = %#v, %v", chapter, err)
	}
}
