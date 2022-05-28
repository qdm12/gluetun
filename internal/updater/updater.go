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

type Updater struct {
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

type Logger interface {
	Info(s string)
	Warn(s string)
	Error(s string)
}

func New(settings settings.Updater, httpClient *http.Client,
	currentServers models.AllServers, logger Logger) *Updater {
	unzipper := unzip.New(httpClient)
	return &Updater{
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

func (u *Updater) UpdateServers(ctx context.Context) (allServers models.AllServers, err error) {
	for _, provider := range u.options.Providers {
		u.logger.Info("updating " + caser.String(provider) + " servers...")
		// TODO support servers offering only TCP or only UDP
		// for NordVPN and PureVPN
		err := u.updateProvider(ctx, provider)
		if err != nil {
			if ctxErr := ctx.Err(); ctxErr != nil {
				return allServers, ctxErr
			}
			u.logger.Error(err.Error())
		}
	}

	return u.servers, nil
}
