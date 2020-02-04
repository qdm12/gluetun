package shadowsocks

import (
	"io"

	"github.com/qdm12/golibs/command"
	"github.com/qdm12/golibs/files"
	"github.com/qdm12/golibs/logging"
)

const logPrefix = "shadowsocks configurator"

type Configurator interface {
	Version() (string, error)
	MakeConf(port uint16, password string) error
	Start(server string, port uint16, password string, log bool) (stdout io.ReadCloser, err error)
}

type configurator struct {
	fileManager files.FileManager
	logger      logging.Logger
	commander   command.Commander
}

func NewConfigurator(fileManager files.FileManager, logger logging.Logger) Configurator {
	return &configurator{fileManager, logger, command.NewCommander()}
}
