package httpproxy

import (
	"context"
	"reflect"

	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/constants"
)

func (l *looper) GetSettings() (settings configuration.HTTPProxy) {
	l.state.settingsMu.RLock()
	defer l.state.settingsMu.RUnlock()
	return l.state.settings
}

func (l *looper) SetSettings(ctx context.Context, settings configuration.HTTPProxy) (
	outcome string) {
	l.state.settingsMu.Lock()
	settingsUnchanged := reflect.DeepEqual(settings, l.state.settings)
	if settingsUnchanged {
		l.state.settingsMu.Unlock()
		return "settings left unchanged"
	}
	newEnabled := settings.Enabled
	previousEnabled := l.state.settings.Enabled
	l.state.settings = settings
	l.state.settingsMu.Unlock()
	// Either restart or set changed status
	switch {
	case !newEnabled && !previousEnabled:
	case newEnabled && previousEnabled:
		_, _ = l.SetStatus(ctx, constants.Stopped)
		_, _ = l.SetStatus(ctx, constants.Running)
	case newEnabled && !previousEnabled:
		_, _ = l.SetStatus(ctx, constants.Running)
	case !newEnabled && previousEnabled:
		_, _ = l.SetStatus(ctx, constants.Stopped)
	}
	return "settings updated"
}
