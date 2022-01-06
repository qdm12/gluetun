package dns

import (
	"context"

	"github.com/qdm12/gluetun/internal/configuration/settings"
)

type SettingsGetSetter interface {
	SettingsGetter
	SettingsSetter
}

type SettingsGetter interface {
	GetSettings() (settings settings.DNS)
}

func (l *Loop) GetSettings() (settings settings.DNS) { return l.state.GetSettings() }

type SettingsSetter interface {
	SetSettings(ctx context.Context, settings settings.DNS) (
		outcome string)
}

func (l *Loop) SetSettings(ctx context.Context, settings settings.DNS) (
	outcome string) {
	return l.state.SetSettings(ctx, settings)
}
