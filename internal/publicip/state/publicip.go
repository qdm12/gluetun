package state

import (
	"net"
)

type PublicIPGetSetter interface {
	PublicIPGetter
	SetPublicIP(publicIP net.IP)
}

type PublicIPGetter interface {
	GetPublicIP() (publicIP net.IP)
}

func (s *State) GetPublicIP() (publicIP net.IP) {
	s.publicIPMu.RLock()
	defer s.publicIPMu.RUnlock()
	publicIP = make(net.IP, len(s.publicIP))
	copy(publicIP, s.publicIP)
	return publicIP
}

func (s *State) SetPublicIP(publicIP net.IP) {
	s.settingsMu.Lock()
	defer s.settingsMu.Unlock()
	s.publicIP = make(net.IP, len(publicIP))
	copy(s.publicIP, publicIP)
}
