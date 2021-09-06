package models

import (
	"net"
)

func (a AllServers) GetCopy() (servers AllServers) {
	servers = a // copy versions and timestamps
	servers.Cyberghost.Servers = a.GetCyberghost()
	servers.Fastestvpn.Servers = a.GetFastestvpn()
	servers.HideMyAss.Servers = a.GetHideMyAss()
	servers.Ipvanish.Servers = a.GetIpvanish()
	servers.Ivpn.Servers = a.GetIvpn()
	servers.Mullvad.Servers = a.GetMullvad()
	servers.Nordvpn.Servers = a.GetNordvpn()
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

func (a *AllServers) GetCyberghost() (servers []CyberghostServer) {
	if a.Cyberghost.Servers == nil {
		return nil
	}
	servers = make([]CyberghostServer, len(a.Cyberghost.Servers))
	for i, serverToCopy := range a.Cyberghost.Servers {
		servers[i] = serverToCopy
		servers[i].IPs = copyIPs(serverToCopy.IPs)
	}
	return servers
}

func (a *AllServers) GetFastestvpn() (servers []FastestvpnServer) {
	if a.Fastestvpn.Servers == nil {
		return nil
	}
	servers = make([]FastestvpnServer, len(a.Fastestvpn.Servers))
	for i, serverToCopy := range a.Fastestvpn.Servers {
		servers[i] = serverToCopy
		servers[i].IPs = copyIPs(serverToCopy.IPs)
	}
	return servers
}

func (a *AllServers) GetHideMyAss() (servers []HideMyAssServer) {
	if a.HideMyAss.Servers == nil {
		return nil
	}
	servers = make([]HideMyAssServer, len(a.HideMyAss.Servers))
	for i, serverToCopy := range a.HideMyAss.Servers {
		servers[i] = serverToCopy
		servers[i].IPs = copyIPs(serverToCopy.IPs)
	}
	return servers
}

func (a *AllServers) GetIpvanish() (servers []IpvanishServer) {
	if a.Ipvanish.Servers == nil {
		return nil
	}
	servers = make([]IpvanishServer, len(a.Ipvanish.Servers))
	for i, serverToCopy := range a.Ipvanish.Servers {
		servers[i] = serverToCopy
		servers[i].IPs = copyIPs(serverToCopy.IPs)
	}
	return servers
}

func (a *AllServers) GetIvpn() (servers []IvpnServer) {
	if a.Ivpn.Servers == nil {
		return nil
	}
	servers = make([]IvpnServer, len(a.Ivpn.Servers))
	for i, serverToCopy := range a.Ivpn.Servers {
		servers[i] = serverToCopy
		servers[i].IPs = copyIPs(serverToCopy.IPs)
	}
	return servers
}

func (a *AllServers) GetMullvad() (servers []MullvadServer) {
	if a.Mullvad.Servers == nil {
		return nil
	}
	servers = make([]MullvadServer, len(a.Mullvad.Servers))
	for i, serverToCopy := range a.Mullvad.Servers {
		servers[i] = serverToCopy
		servers[i].IPs = copyIPs(serverToCopy.IPs)
		servers[i].IPsV6 = copyIPs(serverToCopy.IPsV6)
	}
	return servers
}

func (a *AllServers) GetNordvpn() (servers []NordvpnServer) {
	if a.Nordvpn.Servers == nil {
		return nil
	}
	servers = make([]NordvpnServer, len(a.Nordvpn.Servers))
	for i, serverToCopy := range a.Nordvpn.Servers {
		servers[i] = serverToCopy
		servers[i].IP = copyIP(serverToCopy.IP)
	}
	return servers
}

func (a *AllServers) GetPia() (servers []PIAServer) {
	if a.Pia.Servers == nil {
		return nil
	}
	servers = make([]PIAServer, len(a.Pia.Servers))
	for i, serverToCopy := range a.Pia.Servers {
		servers[i] = serverToCopy
		servers[i].IPs = copyIPs(serverToCopy.IPs)
	}
	return servers
}

func (a *AllServers) GetPrivado() (servers []PrivadoServer) {
	if a.Privado.Servers == nil {
		return nil
	}
	servers = make([]PrivadoServer, len(a.Privado.Servers))
	for i, serverToCopy := range a.Privado.Servers {
		servers[i] = serverToCopy
		servers[i].IP = copyIP(serverToCopy.IP)
	}
	return servers
}

func (a *AllServers) GetPrivatevpn() (servers []PrivatevpnServer) {
	if a.Privatevpn.Servers == nil {
		return nil
	}
	servers = make([]PrivatevpnServer, len(a.Privatevpn.Servers))
	for i, serverToCopy := range a.Privatevpn.Servers {
		servers[i] = serverToCopy
		servers[i].IPs = copyIPs(serverToCopy.IPs)
	}
	return servers
}

func (a *AllServers) GetProtonvpn() (servers []ProtonvpnServer) {
	if a.Protonvpn.Servers == nil {
		return nil
	}
	servers = make([]ProtonvpnServer, len(a.Protonvpn.Servers))
	for i, serverToCopy := range a.Protonvpn.Servers {
		servers[i] = serverToCopy
		servers[i].EntryIP = copyIP(serverToCopy.EntryIP)
		servers[i].ExitIP = copyIP(serverToCopy.ExitIP)
	}
	return servers
}

func (a *AllServers) GetPurevpn() (servers []PurevpnServer) {
	if a.Purevpn.Servers == nil {
		return nil
	}
	servers = make([]PurevpnServer, len(a.Purevpn.Servers))
	for i, serverToCopy := range a.Purevpn.Servers {
		servers[i] = serverToCopy
		servers[i].IPs = copyIPs(serverToCopy.IPs)
	}
	return servers
}

func (a *AllServers) GetSurfshark() (servers []SurfsharkServer) {
	if a.Surfshark.Servers == nil {
		return nil
	}
	servers = make([]SurfsharkServer, len(a.Surfshark.Servers))
	for i, serverToCopy := range a.Surfshark.Servers {
		servers[i] = serverToCopy
		servers[i].IPs = copyIPs(serverToCopy.IPs)
	}
	return servers
}

func (a *AllServers) GetTorguard() (servers []TorguardServer) {
	if a.Torguard.Servers == nil {
		return nil
	}
	servers = make([]TorguardServer, len(a.Torguard.Servers))
	for i, serverToCopy := range a.Torguard.Servers {
		servers[i] = serverToCopy
		servers[i].IPs = copyIPs(serverToCopy.IPs)
	}
	return servers
}

func (a *AllServers) GetVPNUnlimited() (servers []VPNUnlimitedServer) {
	if a.VPNUnlimited.Servers == nil {
		return nil
	}
	servers = make([]VPNUnlimitedServer, len(a.VPNUnlimited.Servers))
	for i, serverToCopy := range a.VPNUnlimited.Servers {
		servers[i] = serverToCopy
		servers[i].IPs = copyIPs(serverToCopy.IPs)
	}
	return servers
}

func (a *AllServers) GetVyprvpn() (servers []VyprvpnServer) {
	if a.Vyprvpn.Servers == nil {
		return nil
	}
	servers = make([]VyprvpnServer, len(a.Vyprvpn.Servers))
	for i, serverToCopy := range a.Vyprvpn.Servers {
		servers[i] = serverToCopy
		servers[i].IPs = copyIPs(serverToCopy.IPs)
	}
	return servers
}

func (a *AllServers) GetWevpn() (servers []WevpnServer) {
	if a.Windscribe.Servers == nil {
		return nil
	}
	servers = make([]WevpnServer, len(a.Wevpn.Servers))
	for i, serverToCopy := range a.Wevpn.Servers {
		servers[i] = serverToCopy
		servers[i].IPs = copyIPs(serverToCopy.IPs)
	}
	return servers
}

func (a *AllServers) GetWindscribe() (servers []WindscribeServer) {
	if a.Windscribe.Servers == nil {
		return nil
	}
	servers = make([]WindscribeServer, len(a.Windscribe.Servers))
	for i, serverToCopy := range a.Windscribe.Servers {
		servers[i] = serverToCopy
		servers[i].IPs = copyIPs(serverToCopy.IPs)
	}
	return servers
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
