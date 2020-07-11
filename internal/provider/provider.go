package provider

import (
	"github.com/qdm12/golibs/files"
	"github.com/qdm12/golibs/network"
	"github.com/qdm12/private-internet-access-docker/internal/constants"
	"github.com/qdm12/private-internet-access-docker/internal/models"
)

// Provider contains methods to read and modify the openvpn configuration to connect as a client
type Provider interface {
	GetOpenVPNConnections(selection models.ServerSelection) (connections []models.OpenVPNConnection, err error)
	BuildConf(connections []models.OpenVPNConnection, verbosity, uid, gid int, root bool, cipher, auth string, extras models.ExtraConfigOptions) (err error)
	GetPortForward() (port uint16, err error)
}

func New(provider models.VPNProvider, client network.Client, fileManager files.FileManager) Provider {
	switch provider {
	case constants.PrivateInternetAccess:
		return newPrivateInternetAccess(client, fileManager)
	case constants.Mullvad:
		return newMullvad(fileManager)
	case constants.Windscribe:
		return newWindscribe(fileManager)
	case constants.Surfshark:
		return newSurfshark(fileManager)
	case constants.Cyberghost:
		return newCyberghost(fileManager)
	default:
		return nil // should never occur
	}
}
