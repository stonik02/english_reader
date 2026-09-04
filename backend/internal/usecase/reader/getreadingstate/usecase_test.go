package getreadingstate

import (
	"context"
	"testing"

	domain "github.com/deniskrylov/english-reader/backend/internal/domain/reader"
	"go.uber.org/mock/gomock"
)

func TestUseCaseFallsBackToFirstChapter(t *testing.T) {
	controller := gomock.NewController(t)
	repository := NewMockReaderRepository(controller)
	repository.EXPECT().EnsureReadyBook(gomock.Any(), "user-1", "book-1").Return(nil)
	repository.EXPECT().Progress(gomock.Any(), "user-1", "book-1").Return(domain.Progress{}, domain.ErrNotFound)
	repository.EXPECT().FirstChapter(gomock.Any(), "book-1").Return(domain.Chapter{ID: "chapter-1", StartCFI: "epubcfi(/6/2)"}, nil)
	repository.EXPECT().Chapter(gomock.Any(), "book-1", "chapter-1", "").Return(domain.Chapter{ID: "chapter-1"}, nil)
	repository.EXPECT().Settings(gomock.Any(), "user-1").Return(domain.Settings{}, nil)
	state, err := New(repository).Execute(context.Background(), "user-1", "book-1")
	if err != nil || state.Progress.ChapterID != "chapter-1" {
		t.Fatalf("Execute() = %#v, %v", state, err)
	}
}

func TestUseCaseRestoresSavedChapter(t *testing.T) {
	controller := gomock.NewController(t)
	repository := NewMockReaderRepository(controller)
	repository.EXPECT().EnsureReadyBook(gomock.Any(), "user-1", "book-1").Return(nil)
	repository.EXPECT().Progress(gomock.Any(), "user-1", "book-1").Return(domain.Progress{ChapterID: "chapter-2", EPUBCFI: "epubcfi(/6/6)"}, nil)
	repository.EXPECT().Chapter(gomock.Any(), "book-1", "chapter-2", "").Return(domain.Chapter{ID: "chapter-2"}, nil)
	repository.EXPECT().Settings(gomock.Any(), "user-1").Return(domain.Settings{}, nil)
	state, err := New(repository).Execute(context.Background(), "user-1", "book-1")
	if err != nil || state.Chapter.ID != "chapter-2" || state.Progress.EPUBCFI != "epubcfi(/6/6)" {
		t.Fatalf("Execute() = %#v, %v", state, err)
	}
}

func TestUseCaseFallsBackWhenReprocessingClearedSavedChapter(t *testing.T) {
	controller := gomock.NewController(t)
	repository := NewMockReaderRepository(controller)
	repository.EXPECT().EnsureReadyBook(gomock.Any(), "user-1", "book-1").Return(nil)
	repository.EXPECT().Progress(gomock.Any(), "user-1", "book-1").Return(domain.Progress{ChapterID: "", EPUBCFI: "epubcfi(/6/10)", ProgressPercent: 25}, nil)
	repository.EXPECT().Chapter(gomock.Any(), "book-1", "", "").Return(domain.Chapter{}, domain.ErrNotFound)
	repository.EXPECT().FirstChapter(gomock.Any(), "book-1").Return(domain.Chapter{ID: "chapter-1", StartCFI: "epubcfi(/6/2)"}, nil)
	repository.EXPECT().Settings(gomock.Any(), "user-1").Return(domain.Settings{}, nil)

	state, err := New(repository).Execute(context.Background(), "user-1", "book-1")
	if err != nil || state.Chapter.ID != "chapter-1" || state.Progress.ChapterID != "chapter-1" {
		t.Fatalf("Execute() = %#v, %v", state, err)
	}
}
