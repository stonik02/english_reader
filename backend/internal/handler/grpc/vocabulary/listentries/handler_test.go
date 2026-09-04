package listentries

import (
	"context"
	"testing"

	readerv1 "github.com/deniskrylov/english-reader/backend/gen/reader/v1"
	domain "github.com/deniskrylov/english-reader/backend/internal/domain/vocabulary"
	uc "github.com/deniskrylov/english-reader/backend/internal/usecase/vocabulary/listentries"
	"go.uber.org/mock/gomock"
)

func TestHandlerListEntriesUsesAuthenticatedSubject(t *testing.T) {
	controller := gomock.NewController(t)
	usecase := NewMockUseCase(controller)
	tokens := NewMockTokenParser(controller)
	tokens.EXPECT().Parse("token").Return("user-1", nil)
	usecase.EXPECT().Execute(gomock.Any(), "user-1", uc.Request{Limit: 20, Query: "go"}).Return(uc.Response{Entries: []domain.Entry{{ID: "entry-1"}}}, nil)

	response, err := New(usecase, tokens).ListEntries(context.Background(), &readerv1.ListEntriesRequest{AccessToken: "token", Limit: 20, Query: "go"})
	if err != nil || len(response.GetEntries()) != 1 || response.GetEntries()[0].GetId() != "entry-1" {
		t.Fatalf("ListEntries() = %#v, %v", response, err)
	}
}
