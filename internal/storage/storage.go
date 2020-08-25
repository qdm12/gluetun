package storage

import (
	"io/ioutil"
	"os"

	"github.com/qdm12/gluetun/internal/models"
)

type Storage interface {
	SyncServers(hardcodedServers models.AllServers) (allServers models.AllServers, err error)
}

type storage struct {
	osStat    func(name string) (os.FileInfo, error)
	readFile  func(filename string) (data []byte, err error)
	writeFile func(filename string, data []byte, perm os.FileMode) error
}

func New() Storage {
	return &storage{
		osStat:    os.Stat,
		readFile:  ioutil.ReadFile,
		writeFile: ioutil.WriteFile,
	}
}
