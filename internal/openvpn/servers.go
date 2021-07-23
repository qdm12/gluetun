package openvpn

import "github.com/qdm12/gluetun/internal/models"

func (l *looper) GetServers() (servers models.AllServers) {
	return l.state.GetServers()
}

func (l *looper) SetServers(servers models.AllServers) {
	l.state.SetServers(servers)
}
