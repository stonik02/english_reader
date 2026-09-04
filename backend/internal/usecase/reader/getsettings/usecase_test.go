package getsettings

import (
	"context"
	"testing"

	domain "github.com/deniskrylov/english-reader/backend/internal/domain/reader"
	"go.uber.org/mock/gomock"
)

func TestUseCaseExecute(t *testing.T) {
	controller := gomock.NewController(t)
	repository := NewMockSettingsRepository(controller)
	repository.EXPECT().Settings(gomock.Any(), "user-1").Return(domain.Settings{FontScale: 100}, nil)
	value, err := New(repository).Execute(context.Background(), "user-1")
	if err != nil || value.FontScale != 100 {
		t.Fatalf("Execute() = %#v, %v", value, err)
	}
}
