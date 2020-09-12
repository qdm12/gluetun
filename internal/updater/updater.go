package updater

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/qdm12/gluetun/internal/models"
)

type Updater interface {
	UpdateServers(ctx context.Context) (allServers models.AllServers, err error)
}

type updater struct {
	// configuration
	options Options

	// state
	servers models.AllServers

	// Functions for tests
	timeNow  func() time.Time
	println  func(s string)
	httpGet  httpGetFunc
	lookupIP lookupIPFunc
}

func New(options Options, httpClient *http.Client, currentServers models.AllServers) Updater {
	if len(options.DNSAddress) == 0 {
		options.DNSAddress = "1.1.1.1"
	}
	resolver := newResolver(options.DNSAddress)
	return &updater{
		timeNow:  time.Now,
		println:  func(s string) { fmt.Println(s) },
		httpGet:  httpClient.Get,
		lookupIP: newLookupIP(resolver),
		options:  options,
		servers:  currentServers,
	}
}

// TODO parallelize DNS resolution
func (u *updater) UpdateServers(ctx context.Context) (allServers models.AllServers, err error) {
	if u.options.Cyberghost {
		u.updateCyberghost(ctx)
	}

	if u.options.Mullvad {
		if err := u.updateMullvad(); err != nil {
			return allServers, err
		}
	}

	if u.options.Nordvpn {
		// TODO support servers offering only TCP or only UDP
		if err := u.updateNordvpn(); err != nil {
			return allServers, err
		}
	}

	if u.options.PIA {
		if err := u.updatePIA(); err != nil {
			return allServers, err
		}
	}

	if u.options.PIAold {
		if err := u.updatePIAOld(ctx); err != nil {
			return allServers, err
		}
	}

	if u.options.Purevpn {
		// TODO support servers offering only TCP or only UDP
		if err := u.updatePurevpn(ctx); err != nil {
			return allServers, err
		}
	}

	if u.options.Surfshark {
		if err := u.updateSurfshark(ctx); err != nil {
			return allServers, err
		}
	}

	if u.options.Vyprvpn {
		if err := u.updateVyprvpn(ctx); err != nil {
			return allServers, err
		}
	}

	if u.options.Windscribe {
		u.updateWindscribe(ctx)
	}

	return u.servers, nil
}
