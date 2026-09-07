package updatesettings

import (
	"context"
	"testing"

	domain "github.com/deniskrylov/english-reader/backend/internal/domain/reader"
	"go.uber.org/mock/gomock"
)

func TestUseCaseRejectsInvalidSettings(t *testing.T) {
	controller := gomock.NewController(t)
	_, err := New(NewMockSettingsRepository(controller)).Execute(context.Background(), "user-1", domain.Settings{FontScale: 20})
	if err != domain.ErrInvalidInput {
		t.Fatalf("error = %v", err)
	}
}

func TestUseCaseExecute(t *testing.T) {
	controller := gomock.NewController(t)
	repository := NewMockSettingsRepository(controller)
	value := domain.Settings{FontScale: 100, Theme: "system", LineHeight: 1.5, HighlightColor: "yellow"}
	repository.EXPECT().UpdateSettings(gomock.Any(), "user-1", value).Return(value, nil)
	if _, err := New(repository).Execute(context.Background(), "user-1", value); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
}

func TestUseCaseAcceptsSubtleAndDisabledHighlights(t *testing.T) {
	for _, color := range []string{"gray-outline", "none"} {
		t.Run(color, func(t *testing.T) {
			controller := gomock.NewController(t)
			repository := NewMockSettingsRepository(controller)
			value := domain.Settings{FontScale: 100, Theme: "system", LineHeight: 1.5, HighlightColor: color}
			repository.EXPECT().UpdateSettings(gomock.Any(), "user-1", value).Return(value, nil)

			_, err := New(repository).Execute(context.Background(), "user-1", value)
			if err != nil {
				t.Fatalf("Execute() error = %v", err)
			}
		})
	}
}

func TestUseCaseRejectsUnknownHighlightColor(t *testing.T) {
	controller := gomock.NewController(t)
	_, err := New(NewMockSettingsRepository(controller)).Execute(context.Background(), "user-1", domain.Settings{FontScale: 100, Theme: "system", LineHeight: 1.5, HighlightColor: "black"})
	if err != domain.ErrInvalidInput {
		t.Fatalf("error = %v", err)
	}
}
