package publicip

import (
	"context"

	"github.com/qdm12/gluetun/internal/configuration/settings"
)

func (l *Loop) GetSettings() (settings settings.PublicIP) {
	return l.state.GetSettings()
}

func (l *Loop) SetSettings(ctx context.Context, settings settings.PublicIP) (
	outcome string) {
	return l.state.SetSettings(ctx, settings)
}
