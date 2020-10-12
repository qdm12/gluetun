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
	GetOpenVPNConnection(selection models.ServerSelection) (connection models.OpenVPNConnection, err error)
	BuildConf(connection models.OpenVPNConnection, verbosity, uid, gid int, root bool, cipher, auth string, extras models.ExtraConfigOptions) (lines []string)
	PortForward(ctx context.Context, client *http.Client,
		fileManager files.FileManager, pfLogger logging.Logger, gateway net.IP, fw firewall.Configurator,
		syncState func(port uint16) (pfFilepath models.Filepath))
}

func New(provider models.VPNProvider, allServers models.AllServers, timeNow timeNowFunc) Provider {
	switch provider {
	case constants.PrivateInternetAccess:
		return newPrivateInternetAccessV4(allServers.Pia.Servers, timeNow)
	case constants.PrivateInternetAccessOld:
		return newPrivateInternetAccessV3(allServers.PiaOld.Servers, timeNow)
	case constants.Mullvad:
		return newMullvad(allServers.Mullvad.Servers, timeNow)
	case constants.Windscribe:
		return newWindscribe(allServers.Windscribe.Servers, timeNow)
	case constants.Surfshark:
		return newSurfshark(allServers.Surfshark.Servers, timeNow)
	case constants.Cyberghost:
		return newCyberghost(allServers.Cyberghost.Servers, timeNow)
	case constants.Vyprvpn:
		return newVyprvpn(allServers.Vyprvpn.Servers, timeNow)
	case constants.Nordvpn:
		return newNordvpn(allServers.Nordvpn.Servers, timeNow)
	case constants.Purevpn:
		return newPurevpn(allServers.Purevpn.Servers, timeNow)
	default:
		return nil // should never occur
	}
}
