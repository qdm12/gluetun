package models

import (
	"net"
)

func (a AllServers) GetCopy() (servers AllServers) {
	servers = a // copy versions and timestamps
	servers.Cyberghost.Servers = a.GetCyberghost()
	servers.Expressvpn.Servers = a.GetExpressvpn()
	servers.Fastestvpn.Servers = a.GetFastestvpn()
	servers.HideMyAss.Servers = a.GetHideMyAss()
	servers.Ipvanish.Servers = a.GetIpvanish()
	servers.Ivpn.Servers = a.GetIvpn()
	servers.Mullvad.Servers = a.GetMullvad()
	servers.Nordvpn.Servers = a.GetNordvpn()
	servers.Perfectprivacy.Servers = a.GetPerfectprivacy()
	servers.Privado.Servers = a.GetPrivado()
	servers.Pia.Servers = a.GetPia()
	servers.Privatevpn.Servers = a.GetPrivatevpn()
	servers.Protonvpn.Servers = a.GetProtonvpn()
	servers.Purevpn.Servers = a.GetPurevpn()
	servers.Surfshark.Servers = a.GetSurfshark()
	servers.Torguard.Servers = a.GetTorguard()
	servers.VPNUnlimited.Servers = a.GetVPNUnlimited()
	servers.Vyprvpn.Servers = a.GetVyprvpn()
	servers.Windscribe.Servers = a.GetWindscribe()
	return servers
}

func (a *AllServers) GetCyberghost() (servers []Server) {
	return copyServers(a.Cyberghost.Servers)
}

func (a *AllServers) GetExpressvpn() (servers []Server) {
	return copyServers(a.Expressvpn.Servers)
}

func (a *AllServers) GetFastestvpn() (servers []Server) {
	return copyServers(a.Fastestvpn.Servers)
}

func (a *AllServers) GetHideMyAss() (servers []Server) {
	return copyServers(a.HideMyAss.Servers)
}

func (a *AllServers) GetIpvanish() (servers []Server) {
	return copyServers(a.Ipvanish.Servers)
}

func (a *AllServers) GetIvpn() (servers []Server) {
	return copyServers(a.Ivpn.Servers)
}

func (a *AllServers) GetMullvad() (servers []Server) {
	return copyServers(a.Mullvad.Servers)
}

func (a *AllServers) GetNordvpn() (servers []Server) {
	return copyServers(a.Nordvpn.Servers)
}

func (a *AllServers) GetPerfectprivacy() (servers []Server) {
	return copyServers(a.Perfectprivacy.Servers)
}

func (a *AllServers) GetPia() (servers []Server) {
	return copyServers(a.Pia.Servers)
}

func (a *AllServers) GetPrivado() (servers []Server) {
	return copyServers(a.Privado.Servers)
}

func (a *AllServers) GetPrivatevpn() (servers []Server) {
	return copyServers(a.Privatevpn.Servers)
}

func (a *AllServers) GetProtonvpn() (servers []Server) {
	return copyServers(a.Protonvpn.Servers)
}

func (a *AllServers) GetPurevpn() (servers []Server) {
	return copyServers(a.Purevpn.Servers)
}

func (a *AllServers) GetSurfshark() (servers []Server) {
	return copyServers(a.Surfshark.Servers)
}

func (a *AllServers) GetTorguard() (servers []Server) {
	return copyServers(a.Torguard.Servers)
}

func (a *AllServers) GetVPNUnlimited() (servers []Server) {
	return copyServers(a.VPNUnlimited.Servers)
}

func (a *AllServers) GetVyprvpn() (servers []Server) {
	return copyServers(a.Vyprvpn.Servers)
}

func (a *AllServers) GetWevpn() (servers []Server) {
	return copyServers(a.Wevpn.Servers)
}

func (a *AllServers) GetWindscribe() (servers []Server) {
	return copyServers(a.Windscribe.Servers)
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
