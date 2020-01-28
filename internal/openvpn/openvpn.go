package openvpn

import (
	"os"

	"github.com/qdm12/golibs/files"
	"github.com/qdm12/golibs/logging"
)

type Configurator interface {
	WriteAuthFile(user, password string) error
	CheckTUN() error
}

type configurator struct {
	fileManager files.FileManager
	logger      logging.Logger
	openFile    func(name string, flag int, perm os.FileMode) (*os.File, error)
}

func NewConfigurator(logger logging.Logger, fileManager files.FileManager) Configurator {
	return &configurator{
		fileManager: fileManager,
		logger:      logger,
		openFile:    os.OpenFile,
	}
}
