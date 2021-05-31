// Package openvpn defines interfaces to interact with openvpn
// and run it in a stateful loop.
package openvpn

import (
	"context"

	"github.com/qdm12/gluetun/internal/unix"
	"github.com/qdm12/golibs/command"
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/os"
)

type Configurator interface {
	Version24(ctx context.Context) (version string, err error)
	Version25(ctx context.Context) (version string, err error)
	WriteAuthFile(user, password string, puid, pgid int) error
	CheckTUN() error
	CreateTUN() error
	Start(ctx context.Context, version string) (
		stdoutLines, stderrLines chan string, waitError chan error, err error)
}

type configurator struct {
	logger    logging.Logger
	commander command.Commander
	os        os.OS
	unix      unix.Unix
}

func NewConfigurator(logger logging.Logger, os os.OS, unix unix.Unix) Configurator {
	return &configurator{
		logger:    logger,
		commander: command.NewCommander(),
		os:        os,
		unix:      unix,
	}
}
