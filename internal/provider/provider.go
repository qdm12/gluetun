// Package provider defines interfaces to interact with each VPN provider.
package provider

import (
	"context"
	"math/rand"
	"net"
	"net/http"
	"time"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/custom"
	"github.com/qdm12/gluetun/internal/provider/cyberghost"
	"github.com/qdm12/gluetun/internal/provider/expressvpn"
	"github.com/qdm12/gluetun/internal/provider/fastestvpn"
	"github.com/qdm12/gluetun/internal/provider/hidemyass"
	"github.com/qdm12/gluetun/internal/provider/ipvanish"
	"github.com/qdm12/gluetun/internal/provider/ivpn"
	"github.com/qdm12/gluetun/internal/provider/mullvad"
	"github.com/qdm12/gluetun/internal/provider/nordvpn"
	"github.com/qdm12/gluetun/internal/provider/perfectprivacy"
	"github.com/qdm12/gluetun/internal/provider/privado"
	"github.com/qdm12/gluetun/internal/provider/privateinternetaccess"
	"github.com/qdm12/gluetun/internal/provider/privatevpn"
	"github.com/qdm12/gluetun/internal/provider/protonvpn"
	"github.com/qdm12/gluetun/internal/provider/purevpn"
	"github.com/qdm12/gluetun/internal/provider/surfshark"
	"github.com/qdm12/gluetun/internal/provider/torguard"
	"github.com/qdm12/gluetun/internal/provider/utils"
	"github.com/qdm12/gluetun/internal/provider/vpnunlimited"
	"github.com/qdm12/gluetun/internal/provider/vyprvpn"
	"github.com/qdm12/gluetun/internal/provider/wevpn"
	"github.com/qdm12/gluetun/internal/provider/windscribe"
)

// Provider contains methods to read and modify the openvpn configuration to connect as a client.
type Provider interface {
	GetConnection(selection settings.ServerSelection) (connection models.Connection, err error)
	BuildConf(connection models.Connection, settings settings.OpenVPN) (lines []string, err error)
	PortForwarder
}

type PortForwarder interface {
	PortForward(ctx context.Context, client *http.Client,
		logger utils.Logger, gateway net.IP, serverName string) (
		port uint16, err error)
	KeepPortForward(ctx context.Context, client *http.Client,
		port uint16, gateway net.IP, serverName string) (err error)
}

func New(provider string, allServers models.AllServers, timeNow func() time.Time) Provider {
	randSource := rand.NewSource(timeNow().UnixNano())
	switch provider {
	case providers.Custom:
		return custom.New()
	case providers.Cyberghost:
		return cyberghost.New(allServers.Cyberghost.Servers, randSource)
	case providers.Expressvpn:
		return expressvpn.New(allServers.Expressvpn.Servers, randSource)
	case providers.Fastestvpn:
		return fastestvpn.New(allServers.Fastestvpn.Servers, randSource)
	case providers.HideMyAss:
		return hidemyass.New(allServers.HideMyAss.Servers, randSource)
	case providers.Ipvanish:
		return ipvanish.New(allServers.Ipvanish.Servers, randSource)
	case providers.Ivpn:
		return ivpn.New(allServers.Ivpn.Servers, randSource)
	case providers.Mullvad:
		return mullvad.New(allServers.Mullvad.Servers, randSource)
	case providers.Nordvpn:
		return nordvpn.New(allServers.Nordvpn.Servers, randSource)
	case providers.Perfectprivacy:
		return perfectprivacy.New(allServers.Perfectprivacy.Servers, randSource)
	case providers.Privado:
		return privado.New(allServers.Privado.Servers, randSource)
	case providers.PrivateInternetAccess:
		return privateinternetaccess.New(allServers.Pia.Servers, randSource, timeNow)
	case providers.Privatevpn:
		return privatevpn.New(allServers.Privatevpn.Servers, randSource)
	case providers.Protonvpn:
		return protonvpn.New(allServers.Protonvpn.Servers, randSource)
	case providers.Purevpn:
		return purevpn.New(allServers.Purevpn.Servers, randSource)
	case providers.Surfshark:
		return surfshark.New(allServers.Surfshark.Servers, randSource)
	case providers.Torguard:
		return torguard.New(allServers.Torguard.Servers, randSource)
	case providers.VPNUnlimited:
		return vpnunlimited.New(allServers.VPNUnlimited.Servers, randSource)
	case providers.Vyprvpn:
		return vyprvpn.New(allServers.Vyprvpn.Servers, randSource)
	case providers.Wevpn:
		return wevpn.New(allServers.Wevpn.Servers, randSource)
	case providers.Windscribe:
		return windscribe.New(allServers.Windscribe.Servers, randSource)
	default:
		panic("provider " + provider + " is unknown") // should never occur
	}
}
