package provider

import (
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/golibs/network"
)

// Provider contains methods to read and modify the openvpn configuration to connect as a client
type Provider interface {
	GetOpenVPNConnections(selection models.ServerSelection) (connections []models.OpenVPNConnection, err error)
	BuildConf(connections []models.OpenVPNConnection, verbosity, uid, gid int, root bool, cipher, auth string, extras models.ExtraConfigOptions) (lines []string)
	GetPortForward(client network.Client) (port uint16, err error)
}

func New(provider models.VPNProvider, allServers models.AllServers) Provider {
	switch provider {
	case constants.PrivateInternetAccess:
		return newPrivateInternetAccess(allServers.Pia.Servers)
	case constants.PrivateInternetAccessOld:
		return newPrivateInternetAccess(allServers.PiaOld.Servers)
	case constants.Mullvad:
		return newMullvad(allServers.Mullvad.Servers)
	case constants.Windscribe:
		return newWindscribe(allServers.Windscribe.Servers)
	case constants.Surfshark:
		return newSurfshark(allServers.Surfshark.Servers)
	case constants.Cyberghost:
		return newCyberghost(allServers.Cyberghost.Servers)
	case constants.Vyprvpn:
		return newVyprvpn(allServers.Vyprvpn.Servers)
	case constants.Nordvpn:
		return newNordvpn(allServers.Nordvpn.Servers)
	case constants.Purevpn:
		return newPurevpn(allServers.Purevpn.Servers)
	default:
		return nil // should never occur
	}
}
