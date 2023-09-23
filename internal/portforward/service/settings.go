package service

import (
	"errors"
	"fmt"
	"net/netip"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/provider"
	"github.com/qdm12/gosettings"
)

type Settings struct {
	UserSettings  settings.PortForwarding
	PortForwarder provider.PortForwarder
	Gateway       netip.Addr // needed for PIA and ProtonVPN
	ServerName    string     // needed for PIA
	Interface     string     // needed for PIA and ProtonVPN, tun0 for example
	VPNProvider   string     // used to validate new settings
}

func (s *Settings) UpdateWith(partialUpdate Settings) (err error) {
	newSettings := s.copy()
	newSettings.overrideWith(partialUpdate)
	err = newSettings.validate()
	if err != nil {
		return fmt.Errorf("validating new settings: %w", err)
	}
	*s = newSettings
	return nil
}

func (s Settings) copy() (copied Settings) {
	copied.UserSettings = s.UserSettings.Copy()
	copied.PortForwarder = s.PortForwarder
	copied.Gateway = s.Gateway
	copied.ServerName = s.ServerName
	copied.Interface = s.Interface
	copied.VPNProvider = s.VPNProvider
	return copied
}

func (s *Settings) overrideWith(update Settings) {
	s.UserSettings.OverrideWith(update.UserSettings)
	s.PortForwarder = gosettings.OverrideWithInterface(s.PortForwarder, update.PortForwarder)
	s.Gateway = gosettings.OverrideWithValidator(s.Gateway, update.Gateway)
	s.ServerName = gosettings.OverrideWithString(s.ServerName, update.ServerName)
	s.Interface = gosettings.OverrideWithString(s.Interface, update.Interface)
	s.VPNProvider = gosettings.OverrideWithString(s.VPNProvider, update.VPNProvider)
}

var ErrVPNProviderNotSet = errors.New("VPN provider not set")

func (s *Settings) validate() (err error) {
	if s.VPNProvider == "" {
		return fmt.Errorf("%w", ErrVPNProviderNotSet)
	}
	return s.UserSettings.Validate(s.VPNProvider)
}
