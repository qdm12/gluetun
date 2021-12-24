package settings

import (
	"fmt"

	"github.com/qdm12/gluetun/internal/models"
)

type Settings struct {
	ControlServer ControlServer `json:"control_server"`
	DNS           DNS           `json:"dns"`
	Firewall      Firewall      `json:"firewall"`
	Health        Health        `json:"health"`
	HTTPProxy     HTTPProxy     `json:"http_proxy"`
	Log           Log           `json:"log"`
	PublicIP      PublicIP      `json:"public_ip"`
	Shadowsocks   Shadowsocks   `json:"shadowsocks"`
	System        System        `json:"system"`
	Updater       Updater       `json:"updater"`
	Version       Version       `json:"version"`
	VPN           VPN           `json:"vpn"`
}

type Source interface {
	Read() (settings Settings, err error)
}

// New populates a settings object using the sources
// starting from the first source and merging in
// unset fields with the next sources.
// It uses the allServers to validate the settings values.
func New(allServers models.AllServers, sources ...Source) (settings Settings, err error) {
	for _, source := range sources {
		settingsFromSource, err := source.Read()
		if err != nil {
			return settings, err
		}
		settings.mergeWith(settingsFromSource)
	}
	settings.setDefaults()

	err = settings.validate(allServers)
	if err != nil {
		return settings, err
	}

	return settings, nil
}

func (s Settings) validate(allServers models.AllServers) (err error) {
	nameToValidation := map[string]func() error{
		"control server":  s.ControlServer.validate,
		"firewall":        s.Firewall.validate,
		"health":          s.Health.validate,
		"http proxy":      s.HTTPProxy.validate,
		"log":             s.Log.validate,
		"public ip check": s.PublicIP.validate,
		"system":          s.System.validate,
		"updater":         s.Updater.validate,
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
		Firewall:      s.Firewall.copy(),
		Health:        s.Health.copy(),
		HTTPProxy:     s.HTTPProxy.copy(),
		Log:           s.Log.copy(),
		PublicIP:      s.PublicIP.copy(),
		System:        s.System.copy(),
		Updater:       s.Updater.copy(),
		Version:       s.Version.copy(),
		VPN:           s.VPN.copy(),
	}
}

func (s *Settings) mergeWith(other Settings) {
	s.ControlServer.mergeWith(other.ControlServer)
	s.Firewall.mergeWith(other.Firewall)
	s.Health.mergeWith(other.Health)
	s.HTTPProxy.mergeWith(other.HTTPProxy)
	s.Log.mergeWith(other.Log)
	s.PublicIP.mergeWith(other.PublicIP)
	s.System.mergeWith(other.System)
	s.Updater.mergeWith(other.Updater)
	s.Version.mergeWith(other.Version)
	s.VPN.mergeWith(other.VPN)
}

func (s *Settings) OverrideWith(other Settings,
	allServers models.AllServers) (err error) {
	patchedSettings := s.copy()
	patchedSettings.ControlServer.overrideWith(other.ControlServer)
	patchedSettings.Firewall.overrideWith(other.Firewall)
	patchedSettings.Health.overrideWith(other.Health)
	patchedSettings.HTTPProxy.overrideWith(other.HTTPProxy)
	patchedSettings.Log.overrideWith(other.Log)
	patchedSettings.PublicIP.overrideWith(other.PublicIP)
	patchedSettings.System.overrideWith(other.System)
	patchedSettings.Updater.overrideWith(other.Updater)
	patchedSettings.Version.overrideWith(other.Version)
	patchedSettings.VPN.overrideWith(other.VPN)
	err = patchedSettings.validate(allServers)
	if err != nil {
		return err
	}
	*s = patchedSettings
	return nil
}

func (s *Settings) setDefaults() {
	s.ControlServer.setDefaults()
	s.Firewall.setDefaults()
	s.Health.setDefaults()
	s.HTTPProxy.setDefaults()
	s.Log.setDefaults()
	s.PublicIP.setDefaults()
	s.System.setDefaults()
	s.Updater.setDefaults()
	s.Version.setDefaults()
	s.VPN.setDefaults()
}
