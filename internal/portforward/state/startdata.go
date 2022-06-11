package state

import (
	"net"

	"github.com/qdm12/gluetun/internal/provider"
)

type StartData struct {
	PortForwarder provider.PortForwarder
	Gateway       net.IP // needed for PIA
	ServerName    string // needed for PIA
	Interface     string // tun0 for example
}

func (s *State) GetStartData() (startData StartData) {
	s.startDataMu.RLock()
	defer s.startDataMu.RUnlock()
	return s.startData
}

func (s *State) SetStartData(startData StartData) {
	s.startDataMu.Lock()
	defer s.startDataMu.Unlock()
	s.startData = startData
}
