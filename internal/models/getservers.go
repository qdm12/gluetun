package models

import (
	"net"
)

func (a AllServers) GetCopy() (allServersCopy AllServers) {
	allServersCopy.Version = a.Version
	allServersCopy.ProviderToServers = make(map[string]Servers, len(a.ProviderToServers))
	for provider, servers := range a.ProviderToServers {
		allServersCopy.ProviderToServers[provider] = Servers{
			Version:   servers.Version,
			Timestamp: servers.Timestamp,
			Servers:   copyServers(servers.Servers),
		}
	}
	return allServersCopy
}

func copyServers(servers []Server) (serversCopy []Server) {
	if servers == nil {
		return nil
	}

	serversCopy = make([]Server, len(servers))
	for i, server := range servers {
		serversCopy[i] = server
		serversCopy[i].IPs = copyIPs(server.IPs)
	}

	return serversCopy
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
