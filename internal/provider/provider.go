// Package provider defines interfaces to interact with each VPN provider.
package provider

import (
	"context"
	"math/rand"
	"net"
	"net/http"
	"time"

	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/custom"
	"github.com/qdm12/gluetun/internal/provider/cyberghost"
	"github.com/qdm12/gluetun/internal/provider/fastestvpn"
	"github.com/qdm12/gluetun/internal/provider/hidemyass"
	"github.com/qdm12/gluetun/internal/provider/ipvanish"
	"github.com/qdm12/gluetun/internal/provider/ivpn"
	"github.com/qdm12/gluetun/internal/provider/mullvad"
	"github.com/qdm12/gluetun/internal/provider/nordvpn"
	"github.com/qdm12/gluetun/internal/provider/privado"
	"github.com/qdm12/gluetun/internal/provider/privateinternetaccess"
	"github.com/qdm12/gluetun/internal/provider/privatevpn"
	"github.com/qdm12/gluetun/internal/provider/protonvpn"
	"github.com/qdm12/gluetun/internal/provider/purevpn"
	"github.com/qdm12/gluetun/internal/provider/surfshark"
	"github.com/qdm12/gluetun/internal/provider/torguard"
	"github.com/qdm12/gluetun/internal/provider/vpnunlimited"
	"github.com/qdm12/gluetun/internal/provider/vyprvpn"
	"github.com/qdm12/gluetun/internal/provider/wevpn"
	"github.com/qdm12/gluetun/internal/provider/windscribe"
	"github.com/qdm12/golibs/logging"
)

// Provider contains methods to read and modify the openvpn configuration to connect as a client.
type Provider interface {
	GetConnection(selection configuration.ServerSelection) (connection models.Connection, err error)
	BuildConf(connection models.Connection, settings configuration.OpenVPN) (lines []string, err error)
	PortForwarder
}

type PortForwarder interface {
	PortForward(ctx context.Context, client *http.Client,
		logger logging.Logger, gateway net.IP, serverName string) (
		port uint16, err error)
	KeepPortForward(ctx context.Context, client *http.Client,
		logger logging.Logger, port uint16, gateway net.IP, serverName string) (
		err error)
}

func New(provider string, allServers models.AllServers, timeNow func() time.Time) Provider {
	randSource := rand.NewSource(timeNow().UnixNano())
	switch provider {
	case constants.Custom:
		return custom.New()
	case constants.Cyberghost:
		return cyberghost.New(allServers.Cyberghost.Servers, randSource)
	case constants.Fastestvpn:
		return fastestvpn.New(allServers.Fastestvpn.Servers, randSource)
	case constants.HideMyAss:
		return hidemyass.New(allServers.HideMyAss.Servers, randSource)
	case constants.Ipvanish:
		return ipvanish.New(allServers.Ipvanish.Servers, randSource)
	case constants.Ivpn:
		return ivpn.New(allServers.Ivpn.Servers, randSource)
	case constants.Mullvad:
		return mullvad.New(allServers.Mullvad.Servers, randSource)
	case constants.Nordvpn:
		return nordvpn.New(allServers.Nordvpn.Servers, randSource)
	case constants.Privado:
		return privado.New(allServers.Privado.Servers, randSource)
	case constants.PrivateInternetAccess:
		return privateinternetaccess.New(allServers.Pia.Servers, randSource, timeNow)
	case constants.Privatevpn:
		return privatevpn.New(allServers.Privatevpn.Servers, randSource)
	case constants.Protonvpn:
		return protonvpn.New(allServers.Protonvpn.Servers, randSource)
	case constants.Purevpn:
		return purevpn.New(allServers.Purevpn.Servers, randSource)
	case constants.Surfshark:
		return surfshark.New(allServers.Surfshark.Servers, randSource)
	case constants.Torguard:
		return torguard.New(allServers.Torguard.Servers, randSource)
	case constants.VPNUnlimited:
		return vpnunlimited.New(allServers.VPNUnlimited.Servers, randSource)
	case constants.Vyprvpn:
		return vyprvpn.New(allServers.Vyprvpn.Servers, randSource)
	case constants.Wevpn:
		return wevpn.New(allServers.Wevpn.Servers, randSource)
	case constants.Windscribe:
		return windscribe.New(allServers.Windscribe.Servers, randSource)
	default:
		return nil // should never occur
	}
}
