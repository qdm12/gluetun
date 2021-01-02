package openvpn

import (
	"context"
	"io"

	"github.com/qdm12/gluetun/internal/unix"
	"github.com/qdm12/golibs/command"
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/os"
)

type Configurator interface {
	Version(ctx context.Context) (string, error)
	WriteAuthFile(user, password string, puid, pgid int) error
	CheckTUN() error
	CreateTUN() error
	Start(ctx context.Context) (stdout io.ReadCloser, waitFn func() error, err error)
}

type configurator struct {
	logger    logging.Logger
	commander command.Commander
	os        os.OS
	unix      unix.Unix
}

func NewConfigurator(logger logging.Logger, os os.OS, unix unix.Unix) Configurator {
	return &configurator{
		logger:    logger.WithPrefix("openvpn configurator: "),
		commander: command.NewCommander(),
		os:        os,
		unix:      unix,
	}
}
