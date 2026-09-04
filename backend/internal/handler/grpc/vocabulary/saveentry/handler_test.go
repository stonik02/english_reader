package saveentry

import (
	"context"
	"testing"

	readerv1 "github.com/deniskrylov/english-reader/backend/gen/reader/v1"
	domain "github.com/deniskrylov/english-reader/backend/internal/domain/vocabulary"
	uc "github.com/deniskrylov/english-reader/backend/internal/usecase/vocabulary/saveentry"
	"go.uber.org/mock/gomock"
)

func TestHandlerSaveEntryUsesAuthenticatedSubject(t *testing.T) {
	controller := gomock.NewController(t)
	usecase := NewMockUseCase(controller)
	tokens := NewMockTokenParser(controller)
	senseID := int64(3)
	tokens.EXPECT().Parse("token").Return("user-1", nil)
	usecase.EXPECT().Execute(gomock.Any(), "user-1", uc.Request{LemmaID: 7, ChosenSenseID: &senseID, SourceForm: "went"}).Return(domain.Entry{ID: "entry-1"}, false, nil)

	response, err := New(usecase, tokens).SaveEntry(context.Background(), &readerv1.SaveEntryRequest{AccessToken: "token", LemmaId: 7, ChosenSenseId: &senseID, SourceForm: "went"})
	if err != nil || response.GetEntry().GetId() != "entry-1" || response.GetAlreadySaved() {
		t.Fatalf("SaveEntry() = %#v, %v", response, err)
	}
}
