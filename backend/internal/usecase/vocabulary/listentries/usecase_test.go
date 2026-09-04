package listentries

import (
	"context"
	"testing"
	"time"

	domain "github.com/deniskrylov/english-reader/backend/internal/domain/vocabulary"
	"go.uber.org/mock/gomock"
)

func TestUseCaseExecuteBuildsNextCursor(t *testing.T) {
	controller := gomock.NewController(t)
	repository := NewMockRepository(controller)
	normalizer := NewMockWordNormalizer(controller)
	normalizer.EXPECT().Normalize("Runs").Return("runs", nil)
	entries := []domain.Entry{
		{ID: "entry-3", CreatedAt: time.Date(2026, 1, 3, 0, 0, 0, 0, time.UTC)},
		{ID: "entry-2", CreatedAt: time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC)},
		{ID: "entry-1", CreatedAt: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)},
	}
	repository.EXPECT().List(gomock.Any(), "user-1", domain.ListRequest{Limit: 3, Query: "runs"}).Return(entries, nil)

	response, err := New(repository, normalizer).Execute(context.Background(), "user-1", Request{Limit: 2, Query: "Runs"})
	if err != nil || len(response.Entries) != 2 || response.NextCursor == "" {
		t.Fatalf("Execute() = %#v, %v", response, err)
	}
}

func TestUseCaseExecuteRejectsMalformedCursor(t *testing.T) {
	controller := gomock.NewController(t)
	response, err := New(NewMockRepository(controller), NewMockWordNormalizer(controller)).Execute(context.Background(), "user-1", Request{Cursor: "%%%"})
	if err != domain.ErrInvalidInput || response.NextCursor != "" {
		t.Fatalf("Execute() = %#v, %v", response, err)
	}
}
