package listcatalog

import (
	"context"
	"testing"
	"time"

	readerv1 "github.com/deniskrylov/english-reader/backend/gen/reader/v1"
	domain "github.com/deniskrylov/english-reader/backend/internal/domain/library"
	"go.uber.org/mock/gomock"
)

func TestHandlerListCatalogReturnsLowercaseBookStatuses(t *testing.T) {
	t.Parallel()

	controller := gomock.NewController(t)
	usecase := NewMockUseCase(controller)
	usecase.EXPECT().Execute(gomock.Any(), "", 20).Return(domain.Page[domain.Book]{
		Items: []domain.Book{
			{ID: "processing", Status: "processing", CreatedAt: time.Now()},
			{ID: "ready", Status: "ready", CreatedAt: time.Now()},
			{ID: "failed", Status: "failed", CreatedAt: time.Now()},
		},
	}, nil)

	response, err := New(usecase).ListCatalog(context.Background(), &readerv1.ListCatalogRequest{Limit: 20})
	if err != nil {
		t.Fatalf("ListCatalog() error = %v", err)
	}

	want := []string{"processing", "ready", "failed"}
	if got := response.GetBooks(); len(got) != len(want) {
		t.Fatalf("ListCatalog() returned %d books, want %d", len(got), len(want))
	} else {
		for index, status := range want {
			if got[index].GetStatus() != status {
				t.Errorf("book %q status = %q, want %q", got[index].GetId(), got[index].GetStatus(), status)
			}
		}
	}
}
