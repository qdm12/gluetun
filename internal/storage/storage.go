// Package storage defines interfaces to interact with the files persisted such as the list of servers.
package storage

import (
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/golibs/logging"
)

type Storage struct {
	mergedServers    models.AllServers
	hardcodedServers models.AllServers
	logger           logging.Logger
	filepath         string
}

// New creates a new storage and reads the servers from the
// embedded servers file and the file on disk.
// Passing an empty filepath disables writing servers to a file.
func New(logger logging.Logger, filepath string) (storage *Storage, err error) {
	// error returned covered by unit test
	harcodedServers, _ := parseHardcodedServers()

	storage = &Storage{
		hardcodedServers: harcodedServers,
		logger:           logger,
		filepath:         filepath,
	}

	if err := storage.SyncServers(); err != nil {
		return nil, err
	}

	return storage, nil
}
