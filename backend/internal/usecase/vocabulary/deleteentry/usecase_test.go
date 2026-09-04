package deleteentry

import (
	"context"
	"testing"

	domain "github.com/deniskrylov/english-reader/backend/internal/domain/vocabulary"
	"go.uber.org/mock/gomock"
)

func TestUseCaseExecuteDeletesByLemma(t *testing.T) {
	controller := gomock.NewController(t)
	repository := NewMockRepository(controller)
	repository.EXPECT().DeleteByLemmaID(gomock.Any(), "user-1", int64(5)).Return(nil)

	if err := New(repository).Execute(context.Background(), "user-1", Request{LemmaID: 5}); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
}

func TestUseCaseExecuteRejectsAmbiguousTarget(t *testing.T) {
	controller := gomock.NewController(t)
	err := New(NewMockRepository(controller)).Execute(context.Background(), "user-1", Request{EntryID: "entry-1", LemmaID: 5})
	if err != domain.ErrInvalidInput {
		t.Fatalf("Execute() error = %v, want %v", err, domain.ErrInvalidInput)
	}
}
