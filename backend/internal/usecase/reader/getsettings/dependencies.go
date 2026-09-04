package getsettings

import (
	"context"
	domain "github.com/deniskrylov/english-reader/backend/internal/domain/reader"
)

type SettingsRepository interface {
	Settings(context.Context, string) (domain.Settings, error)
}
