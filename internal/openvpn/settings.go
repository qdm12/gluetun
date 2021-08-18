package openvpn

import (
	"context"

	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/openvpn/state"
)

type SettingsGetSetter = state.SettingsGetSetter

func (l *Loop) GetSettings() (
	vpn configuration.VPN, provider configuration.Provider) {
	return l.state.GetSettings()
}

func (l *Loop) SetSettings(ctx context.Context,
	vpn configuration.VPN, provider configuration.Provider) (
	outcome string) {
	return l.state.SetSettings(ctx, vpn, provider)
}
