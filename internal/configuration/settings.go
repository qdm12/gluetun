package configuration

import (
	"errors"
	"fmt"
	"strings"

	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/params"
)

// Settings contains all settings for the program to run.
type Settings struct {
	OpenVPN            OpenVPN
	System             System
	DNS                DNS
	Firewall           Firewall
	HTTPProxy          HTTPProxy
	ShadowSocks        ShadowSocks
	Updater            Updater
	PublicIP           PublicIP
	VersionInformation bool
	ControlServer      ControlServer
	Health             Health
}

func (settings *Settings) String() string {
	return strings.Join(settings.lines(), "\n")
}

func (settings *Settings) lines() (lines []string) {
	lines = append(lines, "Settings summary below:")
	lines = append(lines, settings.OpenVPN.lines()...)
	lines = append(lines, settings.DNS.lines()...)
	lines = append(lines, settings.Firewall.lines()...)
	lines = append(lines, settings.System.lines()...)
	lines = append(lines, settings.HTTPProxy.lines()...)
	lines = append(lines, settings.ShadowSocks.lines()...)
	lines = append(lines, settings.Health.lines()...)
	lines = append(lines, settings.ControlServer.lines()...)
	lines = append(lines, settings.Updater.lines()...)
	lines = append(lines, settings.PublicIP.lines()...)
	if settings.VersionInformation {
		lines = append(lines, lastIndent+"Github version information: enabled")
	}
	return lines
}

var (
	ErrOpenvpn       = errors.New("cannot read Openvpn settings")
	ErrSystem        = errors.New("cannot read System settings")
	ErrDNS           = errors.New("cannot read DNS settings")
	ErrFirewall      = errors.New("cannot read firewall settings")
	ErrHTTPProxy     = errors.New("cannot read HTTP proxy settings")
	ErrShadowsocks   = errors.New("cannot read Shadowsocks settings")
	ErrControlServer = errors.New("cannot read control server settings")
	ErrUpdater       = errors.New("cannot read Updater settings")
	ErrPublicIP      = errors.New("cannot read Public IP getter settings")
	ErrHealth        = errors.New("cannot read health settings")
)

// Read obtains all configuration options for the program and returns an error as soon
// as an error is encountered reading them.
func (settings *Settings) Read(env params.Env, logger logging.Logger) (err error) {
	r := newReader(env, logger)

	settings.VersionInformation, err = r.env.OnOff("VERSION_INFORMATION", params.Default("on"))
	if err != nil {
		return fmt.Errorf("environment variable VERSION_INFORMATION: %w", err)
	}

	if err := settings.OpenVPN.read(r); err != nil {
		return fmt.Errorf("%w: %s", ErrOpenvpn, err)
	}

	if err := settings.System.read(r); err != nil {
		return fmt.Errorf("%w: %s", ErrSystem, err)
	}

	if err := settings.DNS.read(r); err != nil {
		return fmt.Errorf("%w: %s", ErrDNS, err)
	}

	if err := settings.Firewall.read(r); err != nil {
		return fmt.Errorf("%w: %s", ErrFirewall, err)
	}

	if err := settings.HTTPProxy.read(r); err != nil {
		return fmt.Errorf("%w: %s", ErrHTTPProxy, err)
	}

	if err := settings.ShadowSocks.read(r); err != nil {
		return fmt.Errorf("%w: %s", ErrShadowsocks, err)
	}

	if err := settings.ControlServer.read(r); err != nil {
		return fmt.Errorf("%w: %s", ErrControlServer, err)
	}

	if err := settings.Updater.read(r); err != nil {
		return fmt.Errorf("%w: %s", ErrUpdater, err)
	}

	if ip := settings.DNS.PlaintextAddress; ip != nil {
		settings.Updater.DNSAddress = ip.String()
	}

	if err := settings.PublicIP.read(r); err != nil {
		return fmt.Errorf("%w: %s", ErrPublicIP, err)
	}

	if err := settings.Health.read(r); err != nil {
		return fmt.Errorf("%w: %s", ErrHealth, err)
	}

	return nil
}
