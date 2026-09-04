package updatesettings

import (
	"context"
	domain "github.com/deniskrylov/english-reader/backend/internal/domain/reader"
)

type SettingsRepository interface {
	UpdateSettings(context.Context, string, domain.Settings) (domain.Settings, error)
}
