// Package storage defines interfaces to interact with the files persisted such as the list of servers.
package storage

import (
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/golibs/logging"
)

type Storage interface {
	// Passing an empty filepath disables writing to a file
	SyncServers(hardcodedServers models.AllServers) (allServers models.AllServers, err error)
	FlushToFile(servers models.AllServers) error
}

type storage struct {
	logger   logging.Logger
	filepath string
}

func New(logger logging.Logger, filepath string) Storage {
	return &storage{
		logger:   logger,
		filepath: filepath,
	}
}
