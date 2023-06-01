package env

import (
	"os"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gosettings/sources/env"
)

type Source struct {
	env    env.Env
	warner Warner
}

type Warner interface {
	Warn(s string)
}

func New(warner Warner) *Source {
	return &Source{
		env:    *env.New(os.Environ()),
		warner: warner,
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

func (s *Source) onRetroActive(oldKey, newKey string) {
	s.warner.Warn(
		"You are using the old environment variable " + oldKey +
			", please consider changing it to " + newKey)
}

// getEnvWithRetro returns the first set environment variable
// key and corresponding value from the environment
// variable keys given. It first goes through the retroKeys
// and end on returning the value corresponding to the currentKey.
// Note retroKeys should be in order from oldest to most
// recent retro-compatibility key.
func (s *Source) getEnvWithRetro(currentKey string,
	retroKeys []string, options ...env.Option) (key string, value *string) {
	// We check retro-compatibility keys first since
	// the current key might be set in the Dockerfile.
	for _, key = range retroKeys {
		value = s.env.Get(key, options...)
		if value != nil {
			s.onRetroActive(key, currentKey)
			return key, value
		}
	}

	return currentKey, s.env.Get(currentKey, options...)
}
