// Package provider defines interfaces to interact with each VPN provider.
package provider

import (
	"context"
	"net"
	"net/http"

	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/firewall"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/os"
)

// Provider contains methods to read and modify the openvpn configuration to connect as a client.
type Provider interface {
	GetOpenVPNConnection(selection configuration.ServerSelection) (connection models.OpenVPNConnection, err error)
	BuildConf(connection models.OpenVPNConnection, username string, settings configuration.OpenVPN) (lines []string)
	PortForward(ctx context.Context, client *http.Client,
		openFile os.OpenFileFunc, pfLogger logging.Logger, gateway net.IP, fw firewall.Configurator,
		syncState func(port uint16) (pfFilepath string))
}

func New(provider string, allServers models.AllServers, timeNow timeNowFunc) Provider {
	switch provider {
	case constants.Cyberghost:
		return newCyberghost(allServers.Cyberghost.Servers, timeNow)
	case constants.HideMyAss:
		return newHideMyAss(allServers.HideMyAss.Servers, timeNow)
	case constants.Mullvad:
		return newMullvad(allServers.Mullvad.Servers, timeNow)
	case constants.Nordvpn:
		return newNordvpn(allServers.Nordvpn.Servers, timeNow)
	case constants.Privado:
		return newPrivado(allServers.Privado.Servers, timeNow)
	case constants.PrivateInternetAccess:
		return newPrivateInternetAccess(allServers.Pia.Servers, timeNow)
	case constants.Privatevpn:
		return newPrivatevpn(allServers.Privatevpn.Servers, timeNow)
	case constants.Purevpn:
		return newPurevpn(allServers.Purevpn.Servers, timeNow)
	case constants.Surfshark:
		return newSurfshark(allServers.Surfshark.Servers, timeNow)
	case constants.Torguard:
		return newTorguard(allServers.Torguard.Servers, timeNow)
	case constants.Vyprvpn:
		return newVyprvpn(allServers.Vyprvpn.Servers, timeNow)
	case constants.Windscribe:
		return newWindscribe(allServers.Windscribe.Servers, timeNow)
	default:
		return nil // should never occur
	}
}
