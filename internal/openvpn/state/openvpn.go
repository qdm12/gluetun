package state

import (
	"context"
	"reflect"

	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/constants"
)

type SettingsGetSetter interface {
	GetSettings() (openvpn configuration.OpenVPN,
		provider configuration.Provider)
	SetSettings(ctx context.Context, openvpn configuration.OpenVPN,
		provider configuration.Provider) (outcome string)
}

func (s *State) GetSettings() (openvpn configuration.OpenVPN,
	provider configuration.Provider) {
	s.settingsMu.RLock()
	openvpn = s.openvpn
	provider = s.provider
	s.settingsMu.RUnlock()
	return openvpn, provider
}

func (s *State) SetSettings(ctx context.Context,
	openvpn configuration.OpenVPN, provider configuration.Provider) (
	outcome string) {
	s.settingsMu.Lock()
	settingsUnchanged := reflect.DeepEqual(s.openvpn, openvpn) &&
		reflect.DeepEqual(s.provider, provider)
	if settingsUnchanged {
		s.settingsMu.Unlock()
		return "settings left unchanged"
	}
	s.openvpn = openvpn
	s.provider = provider
	s.settingsMu.Unlock()
	_, _ = s.statusApplier.ApplyStatus(ctx, constants.Stopped)
	outcome, _ = s.statusApplier.ApplyStatus(ctx, constants.Running)
	return outcome
}
