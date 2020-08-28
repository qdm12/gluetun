package storage

import (
	"io/ioutil"
	"os"

	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/golibs/logging"
)

type Storage interface {
	SyncServers(hardcodedServers models.AllServers, write bool) (allServers models.AllServers, err error)
	FlushToFile(servers models.AllServers) error
}

type storage struct {
	osStat    func(name string) (os.FileInfo, error)
	readFile  func(filename string) (data []byte, err error)
	writeFile func(filename string, data []byte, perm os.FileMode) error
	logger    logging.Logger
}

func New(logger logging.Logger) Storage {
	return &storage{
		osStat:    os.Stat,
		readFile:  ioutil.ReadFile,
		writeFile: ioutil.WriteFile,
		logger:    logger.WithPrefix("storage: "),
	}
}
