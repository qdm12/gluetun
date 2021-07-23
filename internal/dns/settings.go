package dns

import (
	"context"

	"github.com/qdm12/gluetun/internal/configuration"
)

type SettingsGetterSetter interface {
	SettingsGetter
	SettingsSetter
}

type SettingsGetter interface {
	GetSettings() (settings configuration.DNS)
}

func (l *Loop) GetSettings() (settings configuration.DNS) { return l.state.GetSettings() }

type SettingsSetter interface {
	SetSettings(ctx context.Context, settings configuration.DNS) (
		outcome string)
}

func (l *Loop) SetSettings(ctx context.Context, settings configuration.DNS) (
	outcome string) {
	return l.state.SetSettings(ctx, settings)
}
