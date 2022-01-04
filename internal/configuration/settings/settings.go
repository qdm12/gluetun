package settings

import (
	"fmt"

	"github.com/qdm12/gluetun/internal/models"
)

type Settings struct {
	ControlServer ControlServer
	DNS           DNS
	Firewall      Firewall
	Health        Health
	HTTPProxy     HTTPProxy
	Log           Log
	PublicIP      PublicIP
	Shadowsocks   Shadowsocks
	System        System
	Updater       Updater
	Version       Version
	VPN           VPN
}

func (s Settings) Validate(allServers models.AllServers) (err error) {
	nameToValidation := map[string]func() error{
		"control server":  s.ControlServer.validate,
		"dns":             s.DNS.validate,
		"firewall":        s.Firewall.validate,
		"health":          s.Health.Validate,
		"http proxy":      s.HTTPProxy.validate,
		"log":             s.Log.validate,
		"public ip check": s.PublicIP.validate,
		"shadowsocks":     s.Shadowsocks.validate,
		"system":          s.System.validate,
		"updater":         s.Updater.Validate,
		"version":         s.Version.validate,
		"VPN": func() error {
			return s.VPN.validate(allServers)
		},
	}

	for name, validation := range nameToValidation {
		err = validation()
		if err != nil {
			return fmt.Errorf("failed validating %s settings: %w", name, err)
		}
	}

	return nil
}

func (s *Settings) copy() (copied Settings) {
	return Settings{
		ControlServer: s.ControlServer.copy(),
		DNS:           s.DNS.Copy(),
		Firewall:      s.Firewall.copy(),
		Health:        s.Health.copy(),
		HTTPProxy:     s.HTTPProxy.copy(),
		Log:           s.Log.copy(),
		PublicIP:      s.PublicIP.copy(),
		Shadowsocks:   s.Shadowsocks.copy(),
		System:        s.System.copy(),
		Updater:       s.Updater.copy(),
		Version:       s.Version.copy(),
		VPN:           s.VPN.copy(),
	}
}

func (s *Settings) MergeWith(other Settings) {
	s.ControlServer.mergeWith(other.ControlServer)
	s.DNS.mergeWith(other.DNS)
	s.Firewall.mergeWith(other.Firewall)
	s.Health.MergeWith(other.Health)
	s.HTTPProxy.mergeWith(other.HTTPProxy)
	s.Log.mergeWith(other.Log)
	s.PublicIP.mergeWith(other.PublicIP)
	s.Shadowsocks.mergeWith(other.Shadowsocks)
	s.System.mergeWith(other.System)
	s.Updater.mergeWith(other.Updater)
	s.Version.mergeWith(other.Version)
	s.VPN.mergeWith(other.VPN)
}

func (s *Settings) OverrideWith(other Settings,
	allServers models.AllServers) (err error) {
	patchedSettings := s.copy()
	patchedSettings.ControlServer.overrideWith(other.ControlServer)
	patchedSettings.DNS.overrideWith(other.DNS)
	patchedSettings.Firewall.overrideWith(other.Firewall)
	patchedSettings.Health.OverrideWith(other.Health)
	patchedSettings.HTTPProxy.overrideWith(other.HTTPProxy)
	patchedSettings.Log.overrideWith(other.Log)
	patchedSettings.PublicIP.overrideWith(other.PublicIP)
	patchedSettings.Shadowsocks.overrideWith(other.Shadowsocks)
	patchedSettings.System.overrideWith(other.System)
	patchedSettings.Updater.overrideWith(other.Updater)
	patchedSettings.Version.overrideWith(other.Version)
	patchedSettings.VPN.overrideWith(other.VPN)
	err = patchedSettings.Validate(allServers)
	if err != nil {
		return err
	}
	*s = patchedSettings
	return nil
}

func (s *Settings) SetDefaults() {
	s.ControlServer.setDefaults()
	s.DNS.setDefaults()
	s.Firewall.setDefaults()
	s.Health.SetDefaults()
	s.HTTPProxy.setDefaults()
	s.Log.setDefaults()
	s.PublicIP.setDefaults()
	s.Shadowsocks.setDefaults()
	s.System.setDefaults()
	s.Updater.setDefaults()
	s.Version.setDefaults()
	s.VPN.setDefaults()
}
