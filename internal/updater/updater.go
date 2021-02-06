package updater

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/network"
)

type Updater interface {
	UpdateServers(ctx context.Context) (allServers models.AllServers, err error)
}

type updater struct {
	// configuration
	options configuration.Updater

	// state
	servers models.AllServers

	// Functions for tests
	logger   logging.Logger
	timeNow  func() time.Time
	println  func(s string)
	lookupIP lookupIPFunc
	client   network.Client
}

func New(settings configuration.Updater, httpClient *http.Client,
	currentServers models.AllServers, logger logging.Logger) Updater {
	if len(settings.DNSAddress) == 0 {
		settings.DNSAddress = "1.1.1.1"
	}
	resolver := newResolver(settings.DNSAddress)
	const clientTimeout = 10 * time.Second
	return &updater{
		logger:   logger,
		timeNow:  time.Now,
		println:  func(s string) { fmt.Println(s) },
		lookupIP: newLookupIP(resolver),
		client:   network.NewClient(clientTimeout),
		options:  settings,
		servers:  currentServers,
	}
}

// TODO parallelize DNS resolution.
func (u *updater) UpdateServers(ctx context.Context) (allServers models.AllServers, err error) { //nolint:gocognit
	if u.options.Cyberghost {
		u.logger.Info("updating Cyberghost servers...")
		if err := u.updateCyberghost(ctx); err != nil {
			if ctxErr := ctx.Err(); ctxErr != nil {
				return allServers, ctxErr
			}
			u.logger.Error(err)
		}
	}

	if u.options.Mullvad {
		u.logger.Info("updating Mullvad servers...")
		if err := u.updateMullvad(ctx); err != nil {
			u.logger.Error(err)
		}
		if err := ctx.Err(); err != nil {
			return allServers, err
		}
	}

	if u.options.Nordvpn {
		// TODO support servers offering only TCP or only UDP
		u.logger.Info("updating NordVPN servers...")
		if err := u.updateNordvpn(ctx); err != nil {
			u.logger.Error(err)
		}
		if err := ctx.Err(); err != nil {
			return allServers, err
		}
	}

	if u.options.PIA {
		u.logger.Info("updating Private Internet Access servers...")
		if err := u.updatePIA(ctx); err != nil {
			u.logger.Error(err)
		}
		if ctx.Err() != nil {
			return allServers, ctx.Err()
		}
	}

	if u.options.Privado {
		u.logger.Info("updating Privado servers...")
		if err := u.updatePrivado(ctx); err != nil {
			u.logger.Error(err)
		}
		if ctx.Err() != nil {
			return allServers, ctx.Err()
		}
	}

	if u.options.Purevpn {
		u.logger.Info("updating PureVPN servers...")
		// TODO support servers offering only TCP or only UDP
		if err := u.updatePurevpn(ctx); err != nil {
			if ctxErr := ctx.Err(); ctxErr != nil {
				return allServers, ctxErr
			}
			u.logger.Error(err)
		}
	}

	if u.options.Surfshark {
		u.logger.Info("updating Surfshark servers...")
		if err := u.updateSurfshark(ctx); err != nil {
			if ctxErr := ctx.Err(); ctxErr != nil {
				return allServers, ctxErr
			}
			u.logger.Error(err)
		}
	}

	if u.options.Vyprvpn {
		u.logger.Info("updating Vyprvpn servers...")
		if err := u.updateVyprvpn(ctx); err != nil {
			if ctxErr := ctx.Err(); ctxErr != nil {
				return allServers, ctxErr
			}
			u.logger.Error(err)
		}
	}

	if u.options.Windscribe {
		u.logger.Info("updating Windscribe servers...")
		if err := u.updateWindscribe(ctx); err != nil {
			if ctxErr := ctx.Err(); ctxErr != nil {
				return allServers, ctxErr
			}
			u.logger.Error(err)
		}
	}

	return u.servers, nil
}
