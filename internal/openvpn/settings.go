package openvpn

import (
	"context"

	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/openvpn/state"
)

type SettingsGetterSetter = state.SettingsGetterSetter

func (l *Loop) GetSettings() (settings configuration.OpenVPN) {
	return l.state.GetSettings()
}

func (l *Loop) SetSettings(ctx context.Context, settings configuration.OpenVPN) (
	outcome string) {
	return l.state.SetSettings(ctx, settings)
}
