package portforward

import (
	"context"

	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/portforward/state"
)

type SettingsGetSetter = state.SettingsGetSetter

func (l *Loop) GetSettings() (settings configuration.PortForwarding) {
	return l.state.GetSettings()
}

func (l *Loop) SetSettings(ctx context.Context, settings configuration.PortForwarding) (
	outcome string) {
	return l.state.SetSettings(ctx, settings)
}
