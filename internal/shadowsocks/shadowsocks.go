package shadowsocks

import (
	"io"

	"github.com/qdm12/golibs/command"
	"github.com/qdm12/golibs/files"
)

type Configurator interface {
	Version() (string, error)
	MakeConf(port uint16, password string) error
	Start(log bool) (stdout io.ReadCloser, err error)
}

type configurator struct {
	fileManager files.FileManager
	commander   command.Commander
}

func NewConfigurator(fileManager files.FileManager) Configurator {
	return &configurator{fileManager, command.NewCommander()}
}
