// Package openvpn defines interfaces to interact with openvpn
// and run it in a stateful loop.
package openvpn

import (
	"context"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/unix"
	"github.com/qdm12/golibs/command"
	"github.com/qdm12/golibs/logging"
)

type Configurator interface {
	Version24(ctx context.Context) (version string, err error)
	Version25(ctx context.Context) (version string, err error)
	WriteAuthFile(user, password string, puid, pgid int) error
	CheckTUN() error
	CreateTUN() error
	Start(ctx context.Context, version string, flags []string) (
		stdoutLines, stderrLines chan string, waitError chan error, err error)
}

type configurator struct {
	logger       logging.Logger
	cmder        command.RunStarter
	unix         unix.Unix
	authFilePath string
	tunDevPath   string
}

func NewConfigurator(logger logging.Logger, unix unix.Unix,
	cmder command.RunStarter) Configurator {
	return &configurator{
		logger:       logger,
		cmder:        cmder,
		unix:         unix,
		authFilePath: constants.OpenVPNAuthConf,
		tunDevPath:   constants.TunnelDevice,
	}
}
