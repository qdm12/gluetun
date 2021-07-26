// Package openvpn defines interfaces to interact with openvpn
// and run it in a stateful loop.
package openvpn

import (
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/unix"
	"github.com/qdm12/golibs/command"
	"github.com/qdm12/golibs/logging"
)

type Configurator interface {
	VersionGetter
	AuthWriter
	TUNCheckCreater
	Starter
}

type StarterAuthWriter interface {
	Starter
	AuthWriter
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
