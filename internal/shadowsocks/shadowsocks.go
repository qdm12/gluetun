package shadowsocks

import (
	"context"
	"io"

	"github.com/qdm12/golibs/command"
	"github.com/qdm12/golibs/files"
	"github.com/qdm12/golibs/logging"
)

type Configurator interface {
	Version(ctx context.Context) (string, error)
	MakeConf(port uint16, password, method string, uid, gid int) (err error)
	Start(ctx context.Context, server string, port uint16, password string, log bool) (stdout, stderr io.ReadCloser, waitFn func() error, err error)
}

type configurator struct {
	fileManager files.FileManager
	logger      logging.Logger
	commander   command.Commander
}

func NewConfigurator(fileManager files.FileManager, logger logging.Logger) Configurator {
	return &configurator{
		fileManager: fileManager,
		logger:      logger.WithPrefix("shadowsocks configurator: "),
		commander:   command.NewCommander()}
}
