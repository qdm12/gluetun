package openvpn

import (
	"context"

	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/openvpn/state"
)

type SettingsGetterSetter = state.SettingsGetterSetter

func (l *looper) GetSettings() (settings configuration.OpenVPN) {
	return l.state.GetSettings()
}

func (l *looper) SetSettings(ctx context.Context, settings configuration.OpenVPN) (
	outcome string) {
	return l.state.SetSettings(ctx, settings)
}
