// Package updater implements update mechanisms for each VPN provider servers.
package updater

import (
	"context"
	"net/http"
	"time"

	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/updater/resolver"
	"github.com/qdm12/gluetun/internal/updater/unzip"
	"github.com/qdm12/golibs/logging"
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
	logger    logging.Logger
	timeNow   func() time.Time
	presolver resolver.Parallel
	client    *http.Client
	unzipper  unzip.Unzipper
}

func New(settings configuration.Updater, httpClient *http.Client,
	currentServers models.AllServers, logger logging.Logger) Updater {
	if settings.DNSAddress == "" {
		settings.DNSAddress = "1.1.1.1"
	}
	unzipper := unzip.New(httpClient)
	return &updater{
		logger:    logger,
		timeNow:   time.Now,
		presolver: resolver.NewParallelResolver(settings.DNSAddress),
		client:    httpClient,
		unzipper:  unzipper,
		options:   settings,
		servers:   currentServers,
	}
}

//nolint:gocognit,gocyclo
func (u *updater) UpdateServers(ctx context.Context) (allServers models.AllServers, err error) {
	if u.options.Cyberghost {
		u.logger.Info("updating Cyberghost servers...")
		if err := u.updateCyberghost(ctx); err != nil {
			if ctxErr := ctx.Err(); ctxErr != nil {
				return allServers, ctxErr
			}
			u.logger.Error(err.Error())
		}
	}

	if u.options.Fastestvpn {
		u.logger.Info("updating Fastestvpn servers...")
		if err := u.updateFastestvpn(ctx); err != nil {
			u.logger.Error(err.Error())
		}
		if err := ctx.Err(); err != nil {
			return allServers, err
		}
	}

	if u.options.HideMyAss {
		u.logger.Info("updating HideMyAss servers...")
		if err := u.updateHideMyAss(ctx); err != nil {
			u.logger.Error(err.Error())
		}
		if err := ctx.Err(); err != nil {
			return allServers, err
		}
	}

	if u.options.Ipvanish {
		u.logger.Info("updating Ipvanish servers...")
		if err := u.updateIpvanish(ctx); err != nil {
			u.logger.Error(err.Error())
		}
		if err := ctx.Err(); err != nil {
			return allServers, err
		}
	}

	if u.options.Ivpn {
		u.logger.Info("updating Ivpn servers...")
		if err := u.updateIvpn(ctx); err != nil {
			u.logger.Error(err.Error())
		}
		if err := ctx.Err(); err != nil {
			return allServers, err
		}
	}

	if u.options.Mullvad {
		u.logger.Info("updating Mullvad servers...")
		if err := u.updateMullvad(ctx); err != nil {
			u.logger.Error(err.Error())
		}
		if err := ctx.Err(); err != nil {
			return allServers, err
		}
	}

	if u.options.Nordvpn {
		// TODO support servers offering only TCP or only UDP
		u.logger.Info("updating NordVPN servers...")
		if err := u.updateNordvpn(ctx); err != nil {
			u.logger.Error(err.Error())
		}
		if err := ctx.Err(); err != nil {
			return allServers, err
		}
	}

	if u.options.Privado {
		u.logger.Info("updating Privado servers...")
		if err := u.updatePrivado(ctx); err != nil {
			u.logger.Error(err.Error())
		}
		if ctx.Err() != nil {
			return allServers, ctx.Err()
		}
	}

	if u.options.PIA {
		u.logger.Info("updating Private Internet Access servers...")
		if err := u.updatePIA(ctx); err != nil {
			u.logger.Error(err.Error())
		}
		if ctx.Err() != nil {
			return allServers, ctx.Err()
		}
	}

	if u.options.Privatevpn {
		u.logger.Info("updating Privatevpn servers...")
		if err := u.updatePrivatevpn(ctx); err != nil {
			if ctxErr := ctx.Err(); ctxErr != nil {
				return allServers, ctxErr
			}
			u.logger.Error(err.Error())
		}
	}

	if u.options.Protonvpn {
		u.logger.Info("updating Protonvpn servers...")
		if err := u.updateProtonvpn(ctx); err != nil {
			if ctxErr := ctx.Err(); ctxErr != nil {
				return allServers, ctxErr
			}
			u.logger.Error(err.Error())
		}
	}

	if u.options.Purevpn {
		u.logger.Info("updating PureVPN servers...")
		// TODO support servers offering only TCP or only UDP
		if err := u.updatePurevpn(ctx); err != nil {
			if ctxErr := ctx.Err(); ctxErr != nil {
				return allServers, ctxErr
			}
			u.logger.Error(err.Error())
		}
	}

	if u.options.Surfshark {
		u.logger.Info("updating Surfshark servers...")
		if err := u.updateSurfshark(ctx); err != nil {
			if ctxErr := ctx.Err(); ctxErr != nil {
				return allServers, ctxErr
			}
			u.logger.Error(err.Error())
		}
	}

	if u.options.Torguard {
		u.logger.Info("updating Torguard servers...")
		if err := u.updateTorguard(ctx); err != nil {
			if ctxErr := ctx.Err(); ctxErr != nil {
				return allServers, ctxErr
			}
			u.logger.Error(err.Error())
		}
	}

	if u.options.VPNUnlimited {
		u.logger.Info("updating " + constants.VPNUnlimited + " servers...")
		if err := u.updateVPNUnlimited(ctx); err != nil {
			if ctxErr := ctx.Err(); ctxErr != nil {
				return allServers, ctxErr
			}
			u.logger.Error(err.Error())
		}
	}

	if u.options.Vyprvpn {
		u.logger.Info("updating Vyprvpn servers...")
		if err := u.updateVyprvpn(ctx); err != nil {
			if ctxErr := ctx.Err(); ctxErr != nil {
				return allServers, ctxErr
			}
			u.logger.Error(err.Error())
		}
	}

	if u.options.Wevpn {
		u.logger.Info("updating WeVPN servers...")
		if err := u.updateWevpn(ctx); err != nil {
			if ctxErr := ctx.Err(); ctxErr != nil {
				return allServers, ctxErr
			}
			u.logger.Error(err.Error())
		}
	}

	if u.options.Windscribe {
		u.logger.Info("updating Windscribe servers...")
		if err := u.updateWindscribe(ctx); err != nil {
			if ctxErr := ctx.Err(); ctxErr != nil {
				return allServers, ctxErr
			}
			u.logger.Error(err.Error())
		}
	}

	return u.servers, nil
}
