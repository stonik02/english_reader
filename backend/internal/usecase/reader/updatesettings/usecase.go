package updatesettings

import (
	"context"
	domain "github.com/deniskrylov/english-reader/backend/internal/domain/reader"
)

type UseCase struct{ repository SettingsRepository }

func New(repository SettingsRepository) *UseCase { return &UseCase{repository: repository} }
func (u *UseCase) Execute(ctx context.Context, userID string, value domain.Settings) (domain.Settings, error) {
	if value.FontScale < 80 || value.FontScale > 200 || (value.Theme != "system" && value.Theme != "light" && value.Theme != "dark") || value.LineHeight < 1 || value.LineHeight > 3 || !validHighlightColor(value.HighlightColor) {
		return domain.Settings{}, domain.ErrInvalidInput
	}
	return u.repository.UpdateSettings(ctx, userID, value)
}

func validHighlightColor(value string) bool {
	switch value {
	case "yellow", "blue", "green", "pink", "orange", "purple":
		return true
	default:
		return false
	}
}
