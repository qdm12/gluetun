// Package openvpn defines interfaces to interact with openvpn
// and run it in a stateful loop.
package openvpn

import (
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/golibs/command"
	"github.com/qdm12/golibs/logging"
)

type Configurator interface {
	VersionGetter
	AuthWriter
	Starter
}

type StarterAuthWriter interface {
	Starter
	AuthWriter
}

type configurator struct {
	logger       logging.Logger
	cmder        command.RunStarter
	authFilePath string
}

func NewConfigurator(logger logging.Logger,
	cmder command.RunStarter) Configurator {
	return &configurator{
		logger:       logger,
		cmder:        cmder,
		authFilePath: constants.OpenVPNAuthConf,
	}
}
