package state

import (
	"context"
	"reflect"

	"github.com/qdm12/gluetun/internal/configuration/settings"
)

func (s *State) GetSettings() (settings settings.PublicIP) {
	s.settingsMu.RLock()
	defer s.settingsMu.RUnlock()
	return s.settings
}

func (s *State) SetSettings(ctx context.Context, settings settings.PublicIP) (
	outcome string) {
	s.settingsMu.Lock()

	settingsUnchanged := reflect.DeepEqual(s.settings, settings)
	if settingsUnchanged {
		s.settingsMu.Unlock()
		return "settings left unchanged"
	}

	periodChanged := s.settings.Period != settings.Period
	s.settings = settings
	s.settingsMu.Unlock()

	if periodChanged {
		s.updateTicker <- struct{}{}
		// TODO blocking
	}
	return "settings updated"
}
