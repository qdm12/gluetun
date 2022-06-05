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
	storage Storage

	// Functions for tests
	logger    Logger
	timeNow   func() time.Time
	presolver resolver.Parallel
	client    *http.Client
	unzipper  unzip.Unzipper
}

type Storage interface {
	SetServers(provider string, servers []models.Server) (err error)
	GetServersCount(provider string) (count int)
	ServersAreEqual(provider string, servers []models.Server) (equal bool)
}

type Logger interface {
	Info(s string)
	Warn(s string)
	Error(s string)
}

func New(settings settings.Updater, httpClient *http.Client,
	storage Storage, logger Logger) *Updater {
	unzipper := unzip.New(httpClient)
	return &Updater{
		options:   settings,
		storage:   storage,
		logger:    logger,
		timeNow:   time.Now,
		presolver: resolver.NewParallelResolver(settings.DNSAddress.String()),
		client:    httpClient,
		unzipper:  unzipper,
	}
}

func (u *Updater) UpdateServers(ctx context.Context) (err error) {
	caser := cases.Title(language.English)
	for _, provider := range u.options.Providers {
		u.logger.Info("updating " + caser.String(provider) + " servers...")
		// TODO support servers offering only TCP or only UDP
		// for NordVPN and PureVPN
		err := u.updateProvider(ctx, provider)
		if err == nil {
			continue
		}

		// return the only error for the single provider.
		if len(u.options.Providers) == 1 {
			return err
		}

		// stop updating the next providers if context is canceled.
		if ctxErr := ctx.Err(); ctxErr != nil {
			return ctxErr
		}

		// Log the error and continue updating the next provider.
		u.logger.Error(err.Error())
	}

	return nil
}
