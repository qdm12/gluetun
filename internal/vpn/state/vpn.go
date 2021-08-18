package state

import (
	"context"
	"reflect"

	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/constants"
)

type SettingsGetSetter interface {
	GetSettings() (vpn configuration.VPN,
		provider configuration.Provider)
	SetSettings(ctx context.Context, vpn configuration.VPN,
		provider configuration.Provider) (outcome string)
}

func (s *State) GetSettings() (vpn configuration.VPN,
	provider configuration.Provider) {
	s.settingsMu.RLock()
	vpn = s.vpn
	provider = s.provider
	s.settingsMu.RUnlock()
	return vpn, provider
}

func (s *State) SetSettings(ctx context.Context,
	vpn configuration.VPN, provider configuration.Provider) (
	outcome string) {
	s.settingsMu.Lock()
	settingsUnchanged := reflect.DeepEqual(s.vpn, vpn) &&
		reflect.DeepEqual(s.provider, provider)
	if settingsUnchanged {
		s.settingsMu.Unlock()
		return "settings left unchanged"
	}
	s.vpn = vpn
	s.provider = provider
	s.settingsMu.Unlock()
	_, _ = s.statusApplier.ApplyStatus(ctx, constants.Stopped)
	outcome, _ = s.statusApplier.ApplyStatus(ctx, constants.Running)
	return outcome
}
