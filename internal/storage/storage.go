package storage

import (
	"sync"

	"github.com/qdm12/gluetun/internal/models"
)

type Storage struct {
	mergedServers models.AllServers
	mergedMutex   sync.RWMutex
	// this is stored in memory to avoid re-parsing
	// the embedded JSON file on every call to the
	// SyncServers method.
	hardcodedServers models.AllServers
	logger           Infoer
	filepath         string
}

type Infoer interface {
	Info(s string)
}

// New creates a new storage and reads the servers from the
// embedded servers file and the file on disk.
// Passing an empty filepath disables writing servers to a file.
func New(logger Infoer, filepath string) (storage *Storage, err error) {
	// A unit test prevents any error from being returned
	// and ensures all providers are part of the servers returned.
	hardcodedServers, _ := parseHardcodedServers()

	storage = &Storage{
		hardcodedServers: hardcodedServers,
		logger:           logger,
		filepath:         filepath,
	}

	if err := storage.syncServers(); err != nil {
		return nil, err
	}

	return storage, nil
}
