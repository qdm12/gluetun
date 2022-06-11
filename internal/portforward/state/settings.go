package state

import (
	"context"
	"os"
	"reflect"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants"
)

func (s *State) GetSettings() (settings settings.PortForwarding) {
	s.settingsMu.RLock()
	defer s.settingsMu.RUnlock()
	return s.settings
}

func (s *State) SetSettings(ctx context.Context, settings settings.PortForwarding) (
	outcome string) {
	s.settingsMu.Lock()

	settingsUnchanged := reflect.DeepEqual(s.settings, settings)
	if settingsUnchanged {
		s.settingsMu.Unlock()
		return "settings left unchanged"
	}

	if s.settings.Filepath != settings.Filepath {
		_ = os.Rename(*s.settings.Filepath, *settings.Filepath)
	}

	newEnabled := *settings.Enabled
	previousEnabled := *s.settings.Enabled

	s.settings = settings
	s.settingsMu.Unlock()

	switch {
	case !newEnabled && !previousEnabled:
	case newEnabled && previousEnabled:
		// no need to restart for now since we os.Rename the file here.
	case newEnabled && !previousEnabled:
		_, _ = s.statusApplier.ApplyStatus(ctx, constants.Running)
	case !newEnabled && previousEnabled:
		_, _ = s.statusApplier.ApplyStatus(ctx, constants.Stopped)
	}

	return "settings updated"
}
