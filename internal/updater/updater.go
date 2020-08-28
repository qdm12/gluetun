package updater

import (
	"fmt"
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
}

func New(storage storage.Storage) Updater {
	return &updater{
		storage: storage,
		timeNow: time.Now,
		println: func(s string) { fmt.Println(s) },
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

	if options.File {
		if err := u.storage.FlushToFile(allServers); err != nil {
			return fmt.Errorf("cannot update servers: %w", err)
		}
	}

	return nil
}
