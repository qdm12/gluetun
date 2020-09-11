package updater

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/storage"
)

type Updater interface {
	UpdateServers(ctx context.Context) error
}

type updater struct {
	// configuration
	options Options
	storage storage.Storage

	// state
	servers models.AllServers

	// Functions for tests
	timeNow  func() time.Time
	println  func(s string)
	httpGet  httpGetFunc
	lookupIP lookupIPFunc
}

func New(options Options, storage storage.Storage, httpClient *http.Client) Updater {
	if len(options.DNSAddress) == 0 {
		options.DNSAddress = "1.1.1.1"
	}
	resolver := newResolver(options.DNSAddress)
	return &updater{
		storage:  storage,
		timeNow:  time.Now,
		println:  func(s string) { fmt.Println(s) },
		httpGet:  httpClient.Get,
		lookupIP: newLookupIP(resolver),
		options:  options,
	}
}

// TODO parallelize DNS resolution
func (u *updater) UpdateServers(ctx context.Context) (err error) {
	const writeSync = false
	u.servers, err = u.storage.SyncServers(constants.GetAllServers(), writeSync)
	if err != nil {
		return fmt.Errorf("cannot update servers: %w", err)
	}

	if u.options.Cyberghost {
		u.updateCyberghost(ctx)
	}

	if u.options.Mullvad {
		if err := u.updateMullvad(); err != nil {
			return err
		}
	}

	if u.options.Nordvpn {
		// TODO support servers offering only TCP or only UDP
		if err := u.updateNordvpn(); err != nil {
			return err
		}
	}

	if u.options.PIA {
		if err := u.updatePIA(); err != nil {
			return err
		}
	}

	if u.options.PIAold {
		if err := u.updatePIAOld(ctx); err != nil {
			return err
		}
	}

	if u.options.Purevpn {
		// TODO support servers offering only TCP or only UDP
		if err := u.updatePurevpn(ctx); err != nil {
			return err
		}
	}

	if u.options.Surfshark {
		if err := u.updateSurfshark(ctx); err != nil {
			return err
		}
	}

	if u.options.Vyprvpn {
		if err := u.updateVyprvpn(ctx); err != nil {
			return err
		}
	}

	if u.options.Windscribe {
		u.updateWindscribe(ctx)
	}

	if u.options.File {
		if err := u.storage.FlushToFile(u.servers); err != nil {
			return fmt.Errorf("cannot update servers: %w", err)
		}
	}

	return nil
}
