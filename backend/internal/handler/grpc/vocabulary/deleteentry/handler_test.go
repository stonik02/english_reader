package deleteentry

import (
	"context"
	"testing"

	readerv1 "github.com/deniskrylov/english-reader/backend/gen/reader/v1"
	uc "github.com/deniskrylov/english-reader/backend/internal/usecase/vocabulary/deleteentry"
	"go.uber.org/mock/gomock"
)

func TestHandlerDeleteEntryUsesAuthenticatedSubject(t *testing.T) {
	controller := gomock.NewController(t)
	usecase := NewMockUseCase(controller)
	tokens := NewMockTokenParser(controller)
	tokens.EXPECT().Parse("token").Return("user-1", nil)
	usecase.EXPECT().Execute(gomock.Any(), "user-1", uc.Request{LemmaID: 7}).Return(nil)

	response, err := New(usecase, tokens).DeleteEntry(context.Background(), &readerv1.DeleteEntryRequest{AccessToken: "token", Target: &readerv1.DeleteEntryRequest_LemmaId{LemmaId: 7}})
	if err != nil || response == nil {
		t.Fatalf("DeleteEntry() = %#v, %v", response, err)
	}
}
