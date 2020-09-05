package updater

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/storage"
)

type Updater interface {
	UpdateServers(ctx context.Context) error
}

type updater struct {
	options  Options
	storage  storage.Storage
	timeNow  func() time.Time
	println  func(s string)
	httpGet  func(url string) (resp *http.Response, err error)
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
func (u *updater) UpdateServers(ctx context.Context) error {
	const writeSync = false
	allServers, err := u.storage.SyncServers(constants.GetAllServers(), writeSync)
	if err != nil {
		return fmt.Errorf("cannot update servers: %w", err)
	}

	if u.options.PIA {
		const newServers = true
		servers, err := findPIAServers(newServers)
		if err != nil {
			return fmt.Errorf("cannot update PIA servers: %w", err)
		}
		if u.options.Stdout {
			u.println(stringifyPIAServers(servers))
		}
		allServers.Pia.Timestamp = u.timeNow().Unix()
		allServers.Pia.Servers = servers
	}

	if u.options.PIAold {
		const newServers = false
		servers, err := findPIAServers(newServers)
		if err != nil {
			return fmt.Errorf("cannot update PIA old servers: %w", err)
		}
		if u.options.Stdout {
			u.println(stringifyPIAOldServers(servers))
		}
		allServers.PiaOld.Timestamp = u.timeNow().Unix()
		allServers.PiaOld.Servers = servers
	}

	if u.options.Mullvad {
		servers, err := u.findMullvadServers()
		if err != nil {
			return fmt.Errorf("cannot update Mullvad servers: %w", err)
		}
		if u.options.Stdout {
			u.println(stringifyMullvadServers(servers))
		}
		allServers.Mullvad.Timestamp = u.timeNow().Unix()
		allServers.Mullvad.Servers = servers
	}

	if u.options.Vyprvpn {
		servers, err := findVyprvpnServers(ctx, u.lookupIP)
		if err != nil {
			return fmt.Errorf("cannot update Vyprvpn servers: %w", err)
		}
		if u.options.Stdout {
			u.println(stringifyVyprvpnServers(servers))
		}
		allServers.Vyprvpn.Timestamp = u.timeNow().Unix()
		allServers.Vyprvpn.Servers = servers
	}

	if u.options.Surfshark {
		servers, err := findSurfsharkServers(ctx, u.lookupIP)
		if err != nil {
			return fmt.Errorf("cannot update Surfshark servers: %w", err)
		}
		if u.options.Stdout {
			u.println(stringifySurfsharkServers(servers))
		}
		allServers.Surfshark.Timestamp = u.timeNow().Unix()
		allServers.Surfshark.Servers = servers
	}

	if u.options.Nordvpn {
		// TODO support servers offering only TCP or only UDP
		servers, warnings, err := u.findNordvpnServers()
		for _, warning := range warnings {
			u.println(warning)
		}
		if err != nil {
			return fmt.Errorf("cannot update Nordvpn servers: %w", err)
		}
		if u.options.Stdout {
			u.println(stringifyNordvpnServers(servers))
		}
		allServers.Nordvpn.Timestamp = u.timeNow().Unix()
		allServers.Nordvpn.Servers = servers
	}

	if u.options.Purevpn {
		// TODO support servers offering only TCP or only UDP
		servers, warnings, err := u.findPurevpnServers(ctx)
		for _, warning := range warnings {
			u.println(warning)
		}
		if err != nil {
			return fmt.Errorf("cannot update Purevpn servers: %w", err)
		}
		if u.options.Stdout {
			u.println(stringifyPurevpnServers(servers))
		}
		allServers.Purevpn.Timestamp = u.timeNow().Unix()
		allServers.Purevpn.Servers = servers
	}

	if u.options.Windscribe {
		servers := u.findWindscribeServers(ctx)
		if u.options.Stdout {
			u.println(stringifyWindscribeServers(servers))
		}
		allServers.Windscribe.Timestamp = u.timeNow().Unix()
		allServers.Windscribe.Servers = servers
	}

	if u.options.Cyberghost {
		servers := u.findCyberghostServers(ctx)
		if u.options.Stdout {
			u.println(stringifyCyberghostServers(servers))
		}
		allServers.Cyberghost.Timestamp = u.timeNow().Unix()
		allServers.Cyberghost.Servers = servers
	}

	if u.options.File {
		if err := u.storage.FlushToFile(allServers); err != nil {
			return fmt.Errorf("cannot update servers: %w", err)
		}
	}

	return nil
}
