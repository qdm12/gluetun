package updater

import (
	"context"
	"net/http"
	"time"

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
	providers Providers, logger Logger) *Updater {
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
	minRatio float64) (err error) {
	caser := cases.Title(language.English)
	for _, providerName := range providers {
		u.logger.Info("updating " + caser.String(providerName) + " servers...")

		fetcher := u.providers.Get(providerName)
		// TODO support servers offering only TCP or only UDP
		// for NordVPN and PureVPN
		err := u.updateProvider(ctx, fetcher, minRatio)
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
