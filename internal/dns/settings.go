package dns

import (
	"context"

	"github.com/qdm12/gluetun/internal/configuration"
)

func (l *looper) GetSettings() (settings configuration.DNS) { return l.state.GetSettings() }

type SettingsSetter interface {
	SetSettings(ctx context.Context, settings configuration.DNS) (
		outcome string)
}

func (l *looper) SetSettings(ctx context.Context, settings configuration.DNS) (
	outcome string) {
	return l.state.SetSettings(ctx, settings)
}
