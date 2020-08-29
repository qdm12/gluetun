package updater

import (
	"fmt"
	"net/http"
	"time"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/storage"
)

type Updater interface {
	UpdateServers(options Options) error
}

type updater struct {
	storage storage.Storage
	timeNow func() time.Time
	println func(s string)
	httpGet func(url string) (resp *http.Response, err error)
}

func New(storage storage.Storage, httpClient *http.Client) Updater {
	return &updater{
		storage: storage,
		timeNow: time.Now,
		println: func(s string) { fmt.Println(s) },
		httpGet: httpClient.Get,
	}
}

func (u *updater) UpdateServers(options Options) error {
	const writeSync = false
	allServers, err := u.storage.SyncServers(constants.GetAllServers(), writeSync)
	if err != nil {
		return fmt.Errorf("cannot update servers: %w", err)
	}

	if options.PIA {
		const newServers = true
		servers, err := findPIAServers(newServers)
		if err != nil {
			return fmt.Errorf("cannot update PIA servers: %w", err)
		}
		if options.Stdout {
			u.println(stringifyPIAServers(servers))
		}
		allServers.Pia.Timestamp = u.timeNow().Unix()
		allServers.Pia.Servers = servers
	}

	if options.PIAold {
		const newServers = false
		servers, err := findPIAServers(newServers)
		if err != nil {
			return fmt.Errorf("cannot update PIA old servers: %w", err)
		}
		if options.Stdout {
			u.println(stringifyPIAOldServers(servers))
		}
		allServers.PiaOld.Timestamp = u.timeNow().Unix()
		allServers.PiaOld.Servers = servers
	}

	if options.Mullvad {
		servers, err := u.findMullvadServers()
		if err != nil {
			return fmt.Errorf("cannot update Mullvad servers: %w", err)
		}
		if options.Stdout {
			u.println(stringifyMullvadServers(servers))
		}
		allServers.Mullvad.Timestamp = u.timeNow().Unix()
		allServers.Mullvad.Servers = servers
	}

	if options.File {
		if err := u.storage.FlushToFile(allServers); err != nil {
			return fmt.Errorf("cannot update servers: %w", err)
		}
	}

	return nil
}
