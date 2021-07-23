package state

import (
	"context"
	"reflect"

	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/constants"
)

type SettingsGetterSetter interface {
	SettingsGetter
	SettingsSetter
}

type SettingsGetter interface {
	GetSettings() (settings configuration.HTTPProxy)
}

func (s *State) GetSettings() (settings configuration.HTTPProxy) {
	s.settingsMu.RLock()
	defer s.settingsMu.RUnlock()
	return s.settings
}

type SettingsSetter interface {
	SetSettings(ctx context.Context,
		settings configuration.HTTPProxy) (outcome string)
}

func (s *State) SetSettings(ctx context.Context,
	settings configuration.HTTPProxy) (outcome string) {
	s.settingsMu.Lock()
	settingsUnchanged := reflect.DeepEqual(settings, s.settings)
	if settingsUnchanged {
		s.settingsMu.Unlock()
		return "settings left unchanged"
	}
	newEnabled := settings.Enabled
	previousEnabled := s.settings.Enabled
	s.settings = settings
	s.settingsMu.Unlock()
	// Either restart or set changed status
	switch {
	case !newEnabled && !previousEnabled:
	case newEnabled && previousEnabled:
		_, _ = s.statusApplier.ApplyStatus(ctx, constants.Stopped)
		_, _ = s.statusApplier.ApplyStatus(ctx, constants.Running)
	case newEnabled && !previousEnabled:
		_, _ = s.statusApplier.ApplyStatus(ctx, constants.Running)
	case !newEnabled && previousEnabled:
		_, _ = s.statusApplier.ApplyStatus(ctx, constants.Stopped)
	}
	return "settings updated"
}
