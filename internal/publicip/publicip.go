package publicip

import "net"

func (l *Loop) GetPublicIP() (publicIP net.IP) {
	return l.state.GetPublicIP()
}
