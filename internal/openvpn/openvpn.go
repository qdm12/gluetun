package openvpn

import (
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/golibs/command"
)

var _ Interface = (*Configurator)(nil)

type Interface interface {
	VersionGetter
	AuthWriter
	Writer
}

type Configurator struct {
	logger       Infoer
	cmder        command.RunStarter
	configPath   string
	authFilePath string
	puid, pgid   int
}

func New(logger Infoer, cmder command.RunStarter,
	puid, pgid int) *Configurator {
	return &Configurator{
		logger:       logger,
		cmder:        cmder,
		configPath:   configPath,
		authFilePath: constants.OpenVPNAuthConf,
		puid:         puid,
		pgid:         pgid,
	}
}
