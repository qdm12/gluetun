package state

import (
	"context"
	"reflect"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants"
)

func (s *State) GetSettings() (vpn settings.VPN) {
	s.settingsMu.RLock()
	vpn = s.vpn.Copy()
	s.settingsMu.RUnlock()
	return vpn
}

func (s *State) SetSettings(ctx context.Context, vpn settings.VPN) (
	outcome string) {
	s.settingsMu.Lock()
	settingsUnchanged := reflect.DeepEqual(s.vpn, vpn)
	if settingsUnchanged {
		s.settingsMu.Unlock()
		return "settings left unchanged"
	}
	s.vpn = vpn
	s.settingsMu.Unlock()
	_, _ = s.statusApplier.ApplyStatus(ctx, constants.Stopped)
	outcome, _ = s.statusApplier.ApplyStatus(ctx, constants.Running)
	return outcome
}
