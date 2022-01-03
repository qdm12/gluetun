package state

import (
	"context"
	"reflect"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants"
)

type SettingsGetSetter interface {
	GetSettings() (settings settings.DNS)
	SetSettings(ctx context.Context,
		settings settings.DNS) (outcome string)
}

func (s *State) GetSettings() (settings settings.DNS) {
	s.settingsMu.RLock()
	defer s.settingsMu.RUnlock()
	return s.settings
}

func (s *State) SetSettings(ctx context.Context, settings settings.DNS) (
	outcome string) {
	s.settingsMu.Lock()

	settingsUnchanged := reflect.DeepEqual(s.settings, settings)
	if settingsUnchanged {
		s.settingsMu.Unlock()
		return "settings left unchanged"
	}

	// Check for only update period change
	tempSettings := s.settings.Copy()
	*tempSettings.DoT.UpdatePeriod = *settings.DoT.UpdatePeriod
	onlyUpdatePeriodChanged := reflect.DeepEqual(tempSettings, settings)

	s.settings = settings
	s.settingsMu.Unlock()

	if onlyUpdatePeriodChanged {
		s.updateTicker <- struct{}{}
		return "update period changed"
	}

	// Restart
	_, _ = s.statusApplier.ApplyStatus(ctx, constants.Stopped)
	if *settings.DoT.Enabled {
		outcome, _ = s.statusApplier.ApplyStatus(ctx, constants.Running)
	}
	return outcome
}
