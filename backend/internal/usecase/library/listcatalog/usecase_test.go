package listcatalog

import (
	"context"
	"testing"

	domain "github.com/deniskrylov/english-reader/backend/internal/domain/library"
	"go.uber.org/mock/gomock"
)

func TestUseCaseNormalizesLimit(t *testing.T) {
	controller := gomock.NewController(t)
	mock := NewMockBooks(controller)
	mock.EXPECT().ListCatalog(gomock.Any(), "", 100).Return(domain.Page[domain.Book]{}, nil)

	if _, err := New(mock).Execute(context.Background(), "", 1000); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
}
