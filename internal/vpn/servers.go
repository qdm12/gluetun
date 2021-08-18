package vpn

import (
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/vpn/state"
)

type ServersGetterSetter = state.ServersGetterSetter

func (l *Loop) GetServers() (servers models.AllServers) {
	return l.state.GetServers()
}

func (l *Loop) SetServers(servers models.AllServers) {
	l.state.SetServers(servers)
}
