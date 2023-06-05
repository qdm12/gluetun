package env

import (
	"os"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gosettings/sources/env"
)

type Source struct {
	env                 env.Env
	warner              Warner
	handleDeprecatedKey func(deprecatedKey, newKey string)
}

type Warner interface {
	Warn(s string)
}

func New(warner Warner) *Source {
	handleDeprecatedKey := func(deprecatedKey, newKey string) {
		warner.Warn(
			"You are using the old environment variable " + deprecatedKey +
				", please consider changing it to " + newKey)
	}

	return &Source{
		env:                 *env.New(os.Environ(), handleDeprecatedKey),
		warner:              warner,
		handleDeprecatedKey: handleDeprecatedKey,
	}
}

func (s *Source) String() string { return "environment variables" }

func (s *Source) Read() (settings settings.Settings, err error) {
	settings.VPN, err = s.readVPN()
	if err != nil {
		return settings, err
	}

	settings.Firewall, err = s.readFirewall()
	if err != nil {
		return settings, err
	}

	settings.System, err = s.readSystem()
	if err != nil {
		return settings, err
	}

	settings.Health, err = s.ReadHealth()
	if err != nil {
		return settings, err
	}

	settings.HTTPProxy, err = s.readHTTPProxy()
	if err != nil {
		return settings, err
	}

	settings.Log, err = s.readLog()
	if err != nil {
		return settings, err
	}

	settings.PublicIP, err = s.readPublicIP()
	if err != nil {
		return settings, err
	}

	settings.Updater, err = s.readUpdater()
	if err != nil {
		return settings, err
	}

	settings.Version, err = s.readVersion()
	if err != nil {
		return settings, err
	}

	settings.Shadowsocks, err = s.readShadowsocks()
	if err != nil {
		return settings, err
	}

	settings.DNS, err = s.readDNS()
	if err != nil {
		return settings, err
	}

	settings.ControlServer, err = s.readControlServer()
	if err != nil {
		return settings, err
	}

	settings.Pprof, err = s.readPprof()
	if err != nil {
		return settings, err
	}

	return settings, nil
}
