// Package updater implements update mechanisms for each VPN provider servers.
package updater

import (
	"context"
	"net/http"
	"time"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/updater/resolver"
	"github.com/qdm12/gluetun/internal/updater/unzip"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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

var caser = cases.Title(language.English) //nolint:gochecknoglobals

func (u *updater) UpdateServers(ctx context.Context) (allServers models.AllServers, err error) {
	for _, provider := range u.options.Providers {
		u.logger.Info("updating " + caser.String(provider) + " servers...")
		// TODO support servers offering only TCP or only UDP
		// for NordVPN and PureVPN
		warnings, err := u.updateProvider(ctx, provider)
		if *u.options.CLI {
			for _, warning := range warnings {
				u.logger.Warn(provider + ": " + warning)
			}
		}
		if err != nil {
			if ctxErr := ctx.Err(); ctxErr != nil {
				return allServers, ctxErr
			}
			u.logger.Error(err.Error())
		}
	}

	return u.servers, nil
}
