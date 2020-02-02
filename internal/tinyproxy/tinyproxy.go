package tinyproxy

import (
	"io"

	"github.com/qdm12/golibs/command"
	"github.com/qdm12/golibs/files"
	"github.com/qdm12/private-internet-access-docker/internal/models"
)

type Configurator interface {
	Version() (string, error)
	MakeConf(logLevel models.TinyProxyLogLevel, port uint16, user, password string) error
	Start() (stdout io.ReadCloser, err error)
}

type configurator struct {
	commander   command.Commander
	fileManager files.FileManager
}

func NewConfigurator(fileManager files.FileManager) Configurator {
	return &configurator{fileManager: fileManager, commander: command.NewCommander()}
}
