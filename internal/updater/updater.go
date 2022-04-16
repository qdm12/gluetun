// Package updater implements update mechanisms for each VPN provider servers.
package updater

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/updater/resolver"
	"github.com/qdm12/gluetun/internal/updater/unzip"
)

type Updater interface {
	UpdateServers(ctx context.Context) (allServers models.AllServers, err error)
}

type updater struct {
	// configuration
	options settings.Updater

	// state
	servers models.AllServers

	// Functions for tests
	logger    Logger
	timeNow   func() time.Time
	presolver resolver.Parallel
	client    *http.Client
	unzipper  unzip.Unzipper
}

func New(settings settings.Updater, httpClient *http.Client,
	currentServers models.AllServers, logger Logger) Updater {
	unzipper := unzip.New(httpClient)
	return &updater{
		logger:    logger,
		timeNow:   time.Now,
		presolver: resolver.NewParallelResolver(settings.DNSAddress.String()),
		client:    httpClient,
		unzipper:  unzipper,
		options:   settings,
		servers:   currentServers,
	}
}

type updateFunc func(ctx context.Context) (err error)

func (u *updater) UpdateServers(ctx context.Context) (allServers models.AllServers, err error) {
	for _, provider := range u.options.Providers {
		u.logger.Info("updating " + strings.Title(provider) + " servers...")
		updateProvider := u.getUpdateFunction(provider)

		// TODO support servers offering only TCP or only UDP
		// for NordVPN and PureVPN
		err = updateProvider(ctx)
		if err != nil {
			if ctxErr := ctx.Err(); ctxErr != nil {
				return allServers, ctxErr
			}
			u.logger.Error(err.Error())
		}
	}

	return u.servers, nil
}

func (u *updater) getUpdateFunction(provider string) (updateFunction updateFunc) {
	switch provider {
	case providers.Custom:
		panic("cannot update custom provider")
	case providers.Cyberghost:
		return func(ctx context.Context) (err error) { return u.updateCyberghost(ctx) }
	case providers.Expressvpn:
		return func(ctx context.Context) (err error) { return u.updateExpressvpn(ctx) }
	case providers.Fastestvpn:
		return func(ctx context.Context) (err error) { return u.updateFastestvpn(ctx) }
	case providers.HideMyAss:
		return func(ctx context.Context) (err error) { return u.updateHideMyAss(ctx) }
	case providers.Ipvanish:
		return func(ctx context.Context) (err error) { return u.updateIpvanish(ctx) }
	case providers.Ivpn:
		return func(ctx context.Context) (err error) { return u.updateIvpn(ctx) }
	case providers.Mullvad:
		return func(ctx context.Context) (err error) { return u.updateMullvad(ctx) }
	case providers.Nordvpn:
		return func(ctx context.Context) (err error) { return u.updateNordvpn(ctx) }
	case providers.Perfectprivacy:
		return func(ctx context.Context) (err error) { return u.updatePerfectprivacy(ctx) }
	case providers.Privado:
		return func(ctx context.Context) (err error) { return u.updatePrivado(ctx) }
	case providers.PrivateInternetAccess:
		return func(ctx context.Context) (err error) { return u.updatePIA(ctx) }
	case providers.Privatevpn:
		return func(ctx context.Context) (err error) { return u.updatePrivatevpn(ctx) }
	case providers.Protonvpn:
		return func(ctx context.Context) (err error) { return u.updateProtonvpn(ctx) }
	case providers.Purevpn:
		return func(ctx context.Context) (err error) { return u.updatePurevpn(ctx) }
	case providers.Surfshark:
		return func(ctx context.Context) (err error) { return u.updateSurfshark(ctx) }
	case providers.Torguard:
		return func(ctx context.Context) (err error) { return u.updateTorguard(ctx) }
	case providers.VPNUnlimited:
		return func(ctx context.Context) (err error) { return u.updateVPNUnlimited(ctx) }
	case providers.Vyprvpn:
		return func(ctx context.Context) (err error) { return u.updateVyprvpn(ctx) }
	case providers.Wevpn:
		return func(ctx context.Context) (err error) { return u.updateWevpn(ctx) }
	case providers.Windscribe:
		return func(ctx context.Context) (err error) { return u.updateWindscribe(ctx) }
	default:
		panic("provider " + provider + " is unknown")
	}
}
