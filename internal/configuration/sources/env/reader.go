package env

import (
	"os"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/configuration/sources"
)

var _ sources.Source = (*Reader)(nil)

type Reader struct {
	warner Warner
}

type Warner interface {
	Warn(s string)
}

func New(warner Warner) *Reader {
	return &Reader{
		warner: warner,
	}
}

func (r *Reader) Read() (settings settings.Settings, err error) {
	settings.VPN, err = r.readVPN()
	if err != nil {
		return settings, err
	}

	settings.Firewall, err = r.readFirewall()
	if err != nil {
		return settings, err
	}

	settings.System, err = r.readSystem()
	if err != nil {
		return settings, err
	}

	settings.Health, err = r.ReadHealth()
	if err != nil {
		return settings, err
	}

	settings.HTTPProxy, err = r.readHTTPProxy()
	if err != nil {
		return settings, err
	}

	settings.Log, err = readLog()
	if err != nil {
		return settings, err
	}

	settings.PublicIP, err = r.readPublicIP()
	if err != nil {
		return settings, err
	}

	settings.Updater, err = readUpdater()
	if err != nil {
		return settings, err
	}

	settings.Version, err = readVersion()
	if err != nil {
		return settings, err
	}

	settings.Shadowsocks, err = r.readShadowsocks()
	if err != nil {
		return settings, err
	}

	settings.DNS, err = r.readDNS()
	if err != nil {
		return settings, err
	}

	settings.ControlServer, err = r.readControlServer()
	if err != nil {
		return settings, err
	}

	settings.Pprof, err = readPprof()
	if err != nil {
		return settings, err
	}

	return settings, nil
}

func (r *Reader) onRetroActive(oldKey, newKey string) {
	r.warner.Warn(
		"You are using the old environment variable " + oldKey +
			", please consider changing it to " + newKey)
}

// getEnvWithRetro returns the first environment variable
// key and corresponding non empty value from the environment
// variable keys given. It first goes through the retroKeys
// and end on returning the value corresponding to the currentKey.
// Note retroKeys should be in order from oldest to most
// recent retro-compatibility key.
func (r *Reader) getEnvWithRetro(currentKey string,
	retroKeys ...string) (key, value string) {
	// We check retro-compatibility keys first since
	// the current key might be set in the Dockerfile.
	for _, key = range retroKeys {
		value = os.Getenv(key)
		if value != "" {
			r.onRetroActive(key, currentKey)
			return key, value
		}
	}

	return currentKey, os.Getenv(currentKey)
}
