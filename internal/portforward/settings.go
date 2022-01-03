package portforward

import (
	"context"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/portforward/state"
)

type SettingsGetSetter = state.SettingsGetSetter

func (l *Loop) GetSettings() (settings settings.PortForwarding) {
	return l.state.GetSettings()
}

func (l *Loop) SetSettings(ctx context.Context, settings settings.PortForwarding) (
	outcome string) {
	return l.state.SetSettings(ctx, settings)
}
