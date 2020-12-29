package storage

import (
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/os"
	"github.com/qdm12/golibs/logging"
)

type Storage interface {
	SyncServers(hardcodedServers models.AllServers, write bool) (allServers models.AllServers, err error)
	FlushToFile(servers models.AllServers) error
}

type storage struct {
	os     os.OS
	logger logging.Logger
}

func New(logger logging.Logger, os os.OS) Storage {
	return &storage{
		os:     os,
		logger: logger.WithPrefix("storage: "),
	}
}
