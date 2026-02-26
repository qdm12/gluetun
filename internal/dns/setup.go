package dns

import (
	"context"
	"fmt"
	"net/netip"

	"github.com/qdm12/dns/v2/pkg/check"
	"github.com/qdm12/dns/v2/pkg/middlewares/filter/update"
	"github.com/qdm12/dns/v2/pkg/nameserver"
	"github.com/qdm12/dns/v2/pkg/server"
)

func (l *Loop) setupServer(ctx context.Context) (runError <-chan error, err error) {
	settings := l.GetSettings()
	var updateSettings update.Settings
	updateSettings.SetRebindingProtectionExempt(settings.Blacklist.RebindingProtectionExemptHostnames)
	err = l.filter.Update(updateSettings)
	if err != nil {
		return nil, fmt.Errorf("updating filter for rebinding protection: %w", err)
	}

	serverSettings, err := buildServerSettings(settings, l.filter, l.localResolvers, l.logger)
	if err != nil {
		return nil, fmt.Errorf("building server settings: %w", err)
	}

	server, err := server.New(serverSettings)
	if err != nil {
		return nil, fmt.Errorf("creating server: %w", err)
	}

	runError, err = server.Start(ctx)
	if err != nil {
		return nil, fmt.Errorf("starting server: %w", err)
	}
	l.server = server

	// use internal DNS server
	const defaultDNSPort = 53
	nameserver.UseDNSInternally(nameserver.SettingsInternalDNS{
		AddrPort: netip.AddrPortFrom(settings.ServerAddress, defaultDNSPort),
	})
	err = nameserver.UseDNSSystemWide(nameserver.SettingsSystemDNS{
		IPs:        []netip.Addr{settings.ServerAddress},
		ResolvPath: l.resolvConf,
	})
	if err != nil {
		l.logger.Error(err.Error())
	}

	err = check.WaitForDNS(ctx, check.Settings{})
	if err != nil {
		l.stopServer()
		return nil, err
	}

	return runError, nil
}
