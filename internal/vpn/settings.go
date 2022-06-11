package vpn

import (
	"context"

	"github.com/qdm12/gluetun/internal/configuration/settings"
)

func (l *Loop) GetSettings() (settings settings.VPN) {
	return l.state.GetSettings()
}

func (l *Loop) SetSettings(ctx context.Context,
	vpn settings.VPN) (
	outcome string) {
	return l.state.SetSettings(ctx, vpn)
}
