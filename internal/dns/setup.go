package dns

import (
	"context"
	"errors"
	"fmt"

	"github.com/qdm12/dns/v2/pkg/check"
	"github.com/qdm12/dns/v2/pkg/dot"
	"github.com/qdm12/dns/v2/pkg/nameserver"
)

var errUpdateBlockLists = errors.New("cannot update filter block lists")

func (l *Loop) setupUnbound(ctx context.Context) (runError <-chan error, err error) {
	err = l.updateFiles(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errUpdateBlockLists, err)
	}

	settings := l.GetSettings()

	dotSettings, err := buildDoTSettings(settings, l.filter, l.logger)
	if err != nil {
		return nil, fmt.Errorf("building DoT settings: %w", err)
	}

	server, err := dot.NewServer(dotSettings)
	if err != nil {
		return nil, fmt.Errorf("creating DoT server: %w", err)
	}

	runError, err = server.Start()
	if err != nil {
		return nil, fmt.Errorf("starting server: %w", err)
	}
	l.server = server

	// use Unbound
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
		stopErr := l.server.Stop()
		if stopErr != nil {
			l.logger.Error("stopping DoT server: " + stopErr.Error())
		}
		return nil, err
	}

	return runError, nil
}
