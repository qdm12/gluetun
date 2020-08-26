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
	printf  func(format string, a ...interface{}) (n int, err error)
}

func New(storage storage.Storage) Updater {
	return &updater{
		storage: storage,
		timeNow: time.Now,
		printf:  fmt.Printf,
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
			u.printf(stringifyPIAServers(servers))
		}
		allServers.Pia.Timestamp = u.timeNow().Unix()
		allServers.Pia.Servers = servers
	}

	if options.File {
		if err := u.storage.FlushToFile(allServers); err != nil {
			return fmt.Errorf("cannot update servers: %w", err)
		}
	}

	return nil
}
