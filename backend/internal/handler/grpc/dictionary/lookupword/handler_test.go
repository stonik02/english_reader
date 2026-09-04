package lookupword

import (
	"context"
	"testing"

	readerv1 "github.com/deniskrylov/english-reader/backend/gen/reader/v1"
	domain "github.com/deniskrylov/english-reader/backend/internal/domain/dictionary"
	uc "github.com/deniskrylov/english-reader/backend/internal/usecase/dictionary/lookupword"
	"go.uber.org/mock/gomock"
)

func TestHandlerReturnsPartialSuccess(t *testing.T) {
	controller := gomock.NewController(t)
	usecase := NewMockUseCase(controller)
	tokens := NewMockTokenParser(controller)
	tokens.EXPECT().Parse("access-token").Return("user-1", nil)
	usecase.EXPECT().Execute(gomock.Any(), uc.Request{UserID: "user-1", BookID: "book-1", ChapterID: "chapter-1", SelectedText: "went", SentenceText: "We went home."}).Return(domain.LookupResponse{NormalizedLemma: "go", ProviderError: "translation provider unavailable"}, nil)

	response, err := New(usecase, tokens).LookupWord(context.Background(), &readerv1.LookupWordRequest{AccessToken: "access-token", BookId: "book-1", ChapterId: "chapter-1", SelectedText: "went", SentenceText: "We went home."})
	if err != nil {
		t.Fatalf("LookupWord() error = %v", err)
	}
	if response.GetSentenceTranslation().GetProviderError() == "" {
		t.Fatal("provider error was not returned as partial success")
	}
}
