package gethighlights

import (
	"context"
	"testing"

	readerv1 "github.com/deniskrylov/english-reader/backend/gen/reader/v1"
	domain "github.com/deniskrylov/english-reader/backend/internal/domain/vocabulary"
	uc "github.com/deniskrylov/english-reader/backend/internal/usecase/vocabulary/gethighlights"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestHandlerGetHighlights(t *testing.T) {
	controller := gomock.NewController(t)
	usecase := NewMockUseCase(controller)
	tokens := NewMockTokenParser(controller)
	tokens.EXPECT().Parse("access-token").Return("user-1", nil)
	usecase.EXPECT().Execute(gomock.Any(), uc.Request{
		UserID:    "user-1",
		BookID:    "book-1",
		ChapterID: "chapter-1",
	}).Return([]domain.HighlightToken{{
		LemmaID:   11,
		Lemma:     "go",
		Positions: []int{3, 17},
		MatchKind: domain.HighlightMatchKindLemma,
	}}, nil)

	response, err := New(usecase, tokens).GetHighlights(context.Background(), &readerv1.GetHighlightsRequest{
		AccessToken: "access-token",
		BookId:      "book-1",
		ChapterId:   "chapter-1",
	})
	if err != nil {
		t.Fatalf("GetHighlights() error = %v", err)
	}
	if len(response.GetTokens()) != 1 {
		t.Fatalf("GetHighlights() returned %d tokens, want 1", len(response.GetTokens()))
	}
	if response.GetTokens()[0].GetMatchKind() != readerv1.HighlightToken_MATCH_KIND_LEMMA {
		t.Errorf("GetHighlights() match kind = %v, want %v", response.GetTokens()[0].GetMatchKind(), readerv1.HighlightToken_MATCH_KIND_LEMMA)
	}
}

func TestHandlerGetHighlightsMapsUseCaseError(t *testing.T) {
	controller := gomock.NewController(t)
	usecase := NewMockUseCase(controller)
	tokens := NewMockTokenParser(controller)
	tokens.EXPECT().Parse("access-token").Return("user-1", nil)
	usecase.EXPECT().Execute(gomock.Any(), gomock.Any()).Return(nil, domain.ErrNotReady)

	_, err := New(usecase, tokens).GetHighlights(context.Background(), &readerv1.GetHighlightsRequest{AccessToken: "access-token"})
	if status.Code(err) != codes.FailedPrecondition {
		t.Fatalf("GetHighlights() status = %v, want %v", status.Code(err), codes.FailedPrecondition)
	}
}
