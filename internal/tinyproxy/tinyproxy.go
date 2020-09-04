package tinyproxy

import (
	"context"
	"io"

	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/golibs/command"
	"github.com/qdm12/golibs/files"
	"github.com/qdm12/golibs/logging"
)

type Configurator interface {
	Version(ctx context.Context) (string, error)
	MakeConf(logLevel models.TinyProxyLogLevel, port uint16, user, password string, uid, gid int) error
	Start(ctx context.Context) (stdout io.ReadCloser, waitFn func() error, err error)
}

type configurator struct {
	fileManager files.FileManager
	logger      logging.Logger
	commander   command.Commander
}

func NewConfigurator(fileManager files.FileManager, logger logging.Logger) Configurator {
	return &configurator{
		fileManager: fileManager,
		logger:      logger.WithPrefix("tinyproxy configurator: "),
		commander:   command.NewCommander()}
}
