// Package updater implements update mechanisms for each VPN provider servers.
package updater

import (
	"context"
	"net/http"
	"time"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider"
	"github.com/qdm12/gluetun/internal/updater/unzip"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Updater struct {
	// state
	storage Storage

	// Functions for tests
	logger   Logger
	timeNow  func() time.Time
	client   *http.Client
	unzipper unzip.Unzipper
}

type Storage interface {
	SetServers(provider string, servers []models.Server) (err error)
	GetServersCount(provider string) (count int)
	ServersAreEqual(provider string, servers []models.Server) (equal bool)
	// Extra methods to match the provider.New storage interface
	FilterServers(provider string, selection settings.ServerSelection) (filtered []models.Server, err error)
	GetServerByName(provider string, name string) (server models.Server, ok bool)
}

type Logger interface {
	Info(s string)
	Warn(s string)
	Error(s string)
}

func New(httpClient *http.Client,
	storage Storage, logger Logger) *Updater {
	unzipper := unzip.New(httpClient)
	return &Updater{
		storage:  storage,
		logger:   logger,
		timeNow:  time.Now,
		client:   httpClient,
		unzipper: unzipper,
	}
}

func (u *Updater) UpdateServers(ctx context.Context, providers []string) (err error) {
	caser := cases.Title(language.English)
	for _, providerName := range providers {
		u.logger.Info("updating " + caser.String(providerName) + " servers...")

		fetcherStorage := (Storage)(nil) // unused
		fetcher := provider.New(providerName, fetcherStorage, u.timeNow,
			u.logger, u.client, u.unzipper)
		// TODO support servers offering only TCP or only UDP
		// for NordVPN and PureVPN
		err := u.updateProvider(ctx, fetcher)
		if err == nil {
			continue
		}

		// return the only error for the single provider.
		if len(providers) == 1 {
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
