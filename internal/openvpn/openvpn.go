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
	ConfigWriter
}

type configurator struct {
	logger       logging.Logger
	cmder        command.RunStarter
	configPath   string
	authFilePath string
	puid, pgid   int
}

func NewConfigurator(logger logging.Logger,
	cmder command.RunStarter, puid, pgid int) Configurator {
	return &configurator{
		logger:       logger,
		cmder:        cmder,
		configPath:   constants.OpenVPNConf,
		authFilePath: constants.OpenVPNAuthConf,
		puid:         puid,
		pgid:         pgid,
	}
}
