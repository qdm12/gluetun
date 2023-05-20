package storage

import (
	"net/netip"

	"github.com/qdm12/gluetun/internal/models"
)

func copyServer(server models.Server) (serverCopy models.Server) {
	serverCopy = server
	serverCopy.IPs = copyIPs(server.IPs)
	return serverCopy
}

func copyIPs(toCopy []netip.Addr) (copied []netip.Addr) {
	if toCopy == nil {
		return nil
	}

	copied = make([]netip.Addr, len(toCopy))
	copy(copied, toCopy)
	return copied
}
