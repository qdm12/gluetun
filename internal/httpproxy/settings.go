package httpproxy

import (
	"context"

	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/httpproxy/state"
)

type SettingsGetterSetter = state.SettingsGetterSetter

func (l *looper) GetSettings() (settings configuration.HTTPProxy) {
	return l.state.GetSettings()
}

func (l *looper) SetSettings(ctx context.Context, settings configuration.HTTPProxy) (
	outcome string) {
	return l.state.SetSettings(ctx, settings)
}
