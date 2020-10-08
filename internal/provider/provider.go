package provider

import (
	"context"
	"net"
	"net/http"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/firewall"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/golibs/files"
	"github.com/qdm12/golibs/logging"
)

// Provider contains methods to read and modify the openvpn configuration to connect as a client
type Provider interface {
	GetOpenVPNConnections(selection models.ServerSelection) (connections []models.OpenVPNConnection, err error)
	BuildConf(connections []models.OpenVPNConnection, verbosity, uid, gid int, root bool, cipher, auth string, extras models.ExtraConfigOptions) (lines []string)
	PortForward(ctx context.Context, client *http.Client,
		fileManager files.FileManager, pfLogger logging.Logger, gateway net.IP, fw firewall.Configurator,
		syncState func(port uint16) (pfFilepath models.Filepath))
}

func New(provider models.VPNProvider, allServers models.AllServers) Provider {
	switch provider {
	case constants.PrivateInternetAccess:
		return newPrivateInternetAccessV3(allServers.Pia.Servers)
	case constants.PrivateInternetAccessOld:
		return newPrivateInternetAccessV4(allServers.PiaOld.Servers)
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
