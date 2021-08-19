package state

import (
	"context"
	"reflect"

	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/constants"
)

type SettingsGetSetter interface {
	GetSettings() (vpn configuration.VPN)
	SetSettings(ctx context.Context, vpn configuration.VPN) (outcome string)
}

func (s *State) GetSettings() (vpn configuration.VPN) {
	s.settingsMu.RLock()
	vpn = s.vpn
	s.settingsMu.RUnlock()
	return vpn
}

func (s *State) SetSettings(ctx context.Context, vpn configuration.VPN) (
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
