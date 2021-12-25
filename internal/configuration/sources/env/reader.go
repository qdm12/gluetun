package env

import "github.com/qdm12/gluetun/internal/configuration/settings"

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

	settings.Health, err = r.readHealth()
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

	return settings, nil
}

func (r *Reader) onRetroActive(oldKey, newKey string) {
	r.warner.Warn(
		"You are using the old environment variable " + oldKey +
			", please consider changing it to " + newKey)
}
