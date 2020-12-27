package publicip

import "net"

func (l *looper) GetPublicIP() (publicIP net.IP) {
	l.ipMutex.RLock()
	defer l.ipMutex.RUnlock()
	publicIP = make(net.IP, len(l.ip))
	copy(publicIP, l.ip)
	return publicIP
}

func (l *looper) setPublicIP(publicIP net.IP) {
	l.ipMutex.Lock()
	defer l.ipMutex.Unlock()
	l.ip = publicIP
}
