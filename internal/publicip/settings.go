package publicip

import (
	"context"

	"github.com/qdm12/gluetun/internal/configuration"
)

func (l *Loop) GetSettings() (settings configuration.PublicIP) {
	return l.state.GetSettings()
}

func (l *Loop) SetSettings(ctx context.Context, settings configuration.PublicIP) (
	outcome string) {
	return l.state.SetSettings(ctx, settings)
}
