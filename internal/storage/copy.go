package storage

import (
	"net"

	"github.com/qdm12/gluetun/internal/models"
)

func copyServer(server models.Server) (serverCopy models.Server) {
	serverCopy = server
	serverCopy.IPs = copyIPs(server.IPs)
	return serverCopy
}

func copyIPs(toCopy []net.IP) (copied []net.IP) {
	if toCopy == nil {
		return nil
	}

	copied = make([]net.IP, len(toCopy))
	for i := range toCopy {
		copied[i] = copyIP(toCopy[i])
	}

	return copied
}

func copyIP(toCopy net.IP) (copied net.IP) {
	copied = make(net.IP, len(toCopy))
	copy(copied, toCopy)
	return copied
}
