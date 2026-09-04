package saveentry

import (
	"context"
	"testing"

	domain "github.com/deniskrylov/english-reader/backend/internal/domain/vocabulary"
	"go.uber.org/mock/gomock"
)

func TestUseCaseExecuteSavesEntry(t *testing.T) {
	controller := gomock.NewController(t)
	repository := NewMockRepository(controller)
	senseID := int64(12)
	repository.EXPECT().Save(gomock.Any(), "user-1", domain.SaveRequest{
		LemmaID:       7,
		ChosenSenseID: &senseID,
		SourceForm:    "Running",
	}).Return(domain.Entry{ID: "entry-1", LemmaID: 7}, false, nil)

	entry, alreadySaved, err := New(repository).Execute(context.Background(), "user-1", Request{
		LemmaID:       7,
		ChosenSenseID: &senseID,
		SourceForm:    " Running ",
	})
	if err != nil || alreadySaved || entry.ID != "entry-1" {
		t.Fatalf("Execute() = %#v, %t, %v", entry, alreadySaved, err)
	}
}

func TestUseCaseExecuteRejectsInvalidInput(t *testing.T) {
	controller := gomock.NewController(t)
	repository := NewMockRepository(controller)

	_, _, err := New(repository).Execute(context.Background(), "user-1", Request{LemmaID: 0, SourceForm: "word"})
	if err != domain.ErrInvalidInput {
		t.Fatalf("Execute() error = %v, want %v", err, domain.ErrInvalidInput)
	}
}

func TestUseCaseExecuteReturnsDictionaryValidationErrors(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{name: "unknown lemma", err: domain.ErrNotFound},
		{name: "sense from another lemma", err: domain.ErrInvalidSense},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			repository := NewMockRepository(controller)
			repository.EXPECT().Save(gomock.Any(), "user-1", domain.SaveRequest{
				LemmaID:    7,
				SourceForm: "went",
			}).Return(domain.Entry{}, false, test.err)

			_, _, err := New(repository).Execute(context.Background(), "user-1", Request{LemmaID: 7, SourceForm: "went"})
			if err != test.err {
				t.Fatalf("Execute() error = %v, want %v", err, test.err)
			}
		})
	}
}
