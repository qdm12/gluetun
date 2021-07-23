package openvpn

import (
	"context"

	"github.com/qdm12/gluetun/internal/configuration"
)

func (l *looper) GetSettings() (settings configuration.OpenVPN) {
	return l.state.GetSettings()
}

func (l *looper) SetSettings(ctx context.Context, settings configuration.OpenVPN) (
	outcome string) {
	return l.state.SetSettings(ctx, settings)
}
