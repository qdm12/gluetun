package openvpn

import (
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/golibs/command"
	"github.com/qdm12/golibs/logging"
)

var _ Interface = (*Configurator)(nil)

type Interface interface {
	VersionGetter
	AuthWriter
	Runner
	Writer
}

type Configurator struct {
	logger       logging.Logger
	cmder        command.RunStarter
	configPath   string
	authFilePath string
	puid, pgid   int
}

func New(logger logging.Logger,
	cmder command.RunStarter, puid, pgid int) *Configurator {
	return &Configurator{
		logger:       logger,
		cmder:        cmder,
		configPath:   constants.OpenVPNConf,
		authFilePath: constants.OpenVPNAuthConf,
		puid:         puid,
		pgid:         pgid,
	}
}
