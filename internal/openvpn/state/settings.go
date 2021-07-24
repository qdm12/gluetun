package state

import (
	"context"
	"reflect"

	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/constants"
)

type SettingsGetSetter interface {
	GetSettings() (settings configuration.OpenVPN)
	SetSettings(ctx context.Context, settings configuration.OpenVPN) (
		outcome string)
}

func (s *State) GetSettings() (settings configuration.OpenVPN) {
	s.settingsMu.RLock()
	defer s.settingsMu.RUnlock()
	return s.settings
}

func (s *State) SetSettings(ctx context.Context, settings configuration.OpenVPN) (
	outcome string) {
	s.settingsMu.Lock()
	defer s.settingsMu.Unlock()
	settingsUnchanged := reflect.DeepEqual(s.settings, settings)
	if settingsUnchanged {
		return "settings left unchanged"
	}
	s.settings = settings
	_, _ = s.statusApplier.ApplyStatus(ctx, constants.Stopped)
	outcome, _ = s.statusApplier.ApplyStatus(ctx, constants.Running)
	return outcome
}
