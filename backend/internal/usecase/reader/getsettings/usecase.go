package getsettings

import (
	"context"
	domain "github.com/deniskrylov/english-reader/backend/internal/domain/reader"
)

type UseCase struct{ repository SettingsRepository }

func New(repository SettingsRepository) *UseCase { return &UseCase{repository: repository} }
func (u *UseCase) Execute(ctx context.Context, userID string) (domain.Settings, error) {
	return u.repository.Settings(ctx, userID)
}
