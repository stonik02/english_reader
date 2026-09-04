package saveprogress

import (
	"context"
	"testing"

	domain "github.com/deniskrylov/english-reader/backend/internal/domain/reader"
	"go.uber.org/mock/gomock"
)

func TestUseCaseRejectsInvalidCFI(t *testing.T) {
	controller := gomock.NewController(t)
	_, err := New(NewMockReaderRepository(controller)).Execute(context.Background(), "user-1", "book-1", Request{EPUBCFI: "invalid"})
	if err != domain.ErrInvalidInput {
		t.Fatalf("error = %v", err)
	}
}

func TestUseCaseExecute(t *testing.T) {
	controller := gomock.NewController(t)
	repository := NewMockReaderRepository(controller)
	request := Request{ChapterID: "chapter-1", EPUBCFI: "epubcfi(/6/2)", ProgressPercent: 10, Revision: 1}
	repository.EXPECT().EnsureReadyBook(gomock.Any(), "user-1", "book-1").Return(nil)
	repository.EXPECT().Chapter(gomock.Any(), "book-1", "chapter-1", "").Return(domain.Chapter{}, nil)
	repository.EXPECT().SaveProgress(gomock.Any(), "user-1", "book-1", gomock.Any()).Return(domain.Progress{Revision: 1}, nil)
	if _, err := New(repository).Execute(context.Background(), "user-1", "book-1", request); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
}
