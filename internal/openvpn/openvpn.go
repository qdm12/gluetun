package openvpn

import (
	"io"
	"os"

	"github.com/qdm12/golibs/command"
	"github.com/qdm12/golibs/files"
	"github.com/qdm12/golibs/logging"
)

const logPrefix = "openvpn configurator"

type Configurator interface {
	Version() (string, error)
	WriteAuthFile(user, password string, uid, gid int) error
	CheckTUN() error
	Start() (stdout io.ReadCloser, err error)
}

type configurator struct {
	fileManager files.FileManager
	logger      logging.Logger
	commander   command.Commander
	openFile    func(name string, flag int, perm os.FileMode) (*os.File, error)
}

func NewConfigurator(logger logging.Logger, fileManager files.FileManager) Configurator {
	return &configurator{
		fileManager: fileManager,
		logger:      logger,
		commander:   command.NewCommander(),
		openFile:    os.OpenFile,
	}
}
