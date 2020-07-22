package provider

import (
	"github.com/qdm12/golibs/network"
	"github.com/qdm12/private-internet-access-docker/internal/constants"
	"github.com/qdm12/private-internet-access-docker/internal/models"
)

// Provider contains methods to read and modify the openvpn configuration to connect as a client
type Provider interface {
	GetOpenVPNConnections(selection models.ServerSelection) (connections []models.OpenVPNConnection, err error)
	BuildConf(connections []models.OpenVPNConnection, verbosity, uid, gid int, root bool, cipher, auth string, extras models.ExtraConfigOptions) (lines []string)
	GetPortForward(client network.Client) (port uint16, err error)
}

func New(provider models.VPNProvider) Provider {
	switch provider {
	case constants.PrivateInternetAccess:
		return newPrivateInternetAccess()
	case constants.Mullvad:
		return newMullvad()
	case constants.Windscribe:
		return newWindscribe()
	case constants.Surfshark:
		return newSurfshark()
	case constants.Cyberghost:
		return newCyberghost()
	case constants.Vyprvpn:
		return newVyprvpn()
	case constants.Nordvpn:
		return newNordvpn()
	case constants.Purevpn:
		return newPurevpn()
	default:
		return nil // should never occur
	}
}
