package openvpn

import (
	"github.com/qdm12/gluetun/internal/constants/openvpn"
	"github.com/qdm12/golibs/command"
)

type Configurator struct {
	logger       Infoer
	cmder        command.RunStarter
	configPath   string
	authFilePath string
	askPassPath  string
	puid, pgid   int
}

func New(logger Infoer, cmder command.RunStarter,
	puid, pgid int) *Configurator {
	return &Configurator{
		logger:       logger,
		cmder:        cmder,
		configPath:   configPath,
		authFilePath: openvpn.AuthConf,
		askPassPath:  openvpn.AskPassPath,
		puid:         puid,
		pgid:         pgid,
	}
}
