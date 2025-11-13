package updater

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/qdm12/gluetun/internal/provider/common"
	"github.com/qdm12/gluetun/internal/updater/unzip"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Updater struct {
	providers Providers

	// state
	storage Storage

	// Functions for tests
	logger   Logger
	timeNow  func() time.Time
	client   *http.Client
	unzipper Unzipper
}

func New(httpClient *http.Client, storage Storage,
	providers Providers, logger Logger,
) *Updater {
	unzipper := unzip.New(httpClient)
	return &Updater{
		providers: providers,
		storage:   storage,
		logger:    logger,
		timeNow:   time.Now,
		client:    httpClient,
		unzipper:  unzipper,
	}
}

func (u *Updater) UpdateServers(ctx context.Context, providers []string,
	minRatio float64,
) (err error) {
	caser := cases.Title(language.English)
	for _, providerName := range providers {
		u.logger.Info("updating " + caser.String(providerName) + " servers...")

		fetcher := u.providers.Get(providerName)
		// TODO support servers offering only TCP or only UDP
		// for NordVPN and PureVPN
		err := u.updateProvider(ctx, fetcher, minRatio)
		switch {
		case err == nil:
			continue
		case errors.Is(err, common.ErrCredentialsMissing):
			u.logger.Warn(err.Error() + " - skipping update for " + providerName)
			continue
		case len(providers) == 1:
			// return the only error for the single provider.
			return err
		case ctx.Err() != nil:
			// stop updating other providers if context is done
			return ctx.Err()
		default: // error encountered updating one of multiple providers
			// Log the error and continue updating the next provider.
			u.logger.Error(err.Error())
		}
	}

	return nil
}
