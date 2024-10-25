package dns

import (
	"context"
	"errors"
	"fmt"

	"github.com/qdm12/dns/v2/pkg/check"
	"github.com/qdm12/dns/v2/pkg/nameserver"
	"github.com/qdm12/dns/v2/pkg/server"
)

var errUpdateBlockLists = errors.New("cannot update filter block lists")

func (l *Loop) setupServer(ctx context.Context) (runError <-chan error, err error) {
	err = l.updateFiles(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errUpdateBlockLists, err)
	}

	settings := l.GetSettings()

	dotSettings, err := buildDoTSettings(settings, l.filter, l.logger)
	if err != nil {
		return nil, fmt.Errorf("building DoT settings: %w", err)
	}

	server, err := server.New(dotSettings)
	if err != nil {
		return nil, fmt.Errorf("creating DoT server: %w", err)
	}

	runError, err = server.Start(ctx)
	if err != nil {
		return nil, fmt.Errorf("starting server: %w", err)
	}
	l.server = server

	// use internal DNS server
	nameserver.UseDNSInternally(nameserver.SettingsInternalDNS{
		IP: settings.ServerAddress,
	})
	err = nameserver.UseDNSSystemWide(nameserver.SettingsSystemDNS{
		IP:         settings.ServerAddress,
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
