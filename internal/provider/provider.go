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
	"github.com/qdm12/gluetun/internal/provider/common"
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
	OpenVPNConfig(connection models.Connection, settings settings.OpenVPN) (lines []string)
	Name() string
	PortForwarder
	FetchServers(ctx context.Context, minServers int) (
		servers []models.Server, err error)
}

type PortForwarder interface {
	PortForward(ctx context.Context, client *http.Client,
		logger utils.Logger, gateway net.IP, serverName string) (
		port uint16, err error)
	KeepPortForward(ctx context.Context, client *http.Client,
		port uint16, gateway net.IP, serverName string) (err error)
}

type Storage interface {
	FilterServers(provider string, selection settings.ServerSelection) (
		servers []models.Server, err error)
	GetServerByName(provider, name string) (server models.Server, ok bool)
}

func New(provider string, storage Storage, timeNow func() time.Time,
	updaterWarner common.Warner, client *http.Client, unzipper common.Unzipper) Provider {
	randSource := rand.NewSource(timeNow().UnixNano())
	switch provider {
	case providers.Custom:
		return custom.New()
	case providers.Cyberghost:
		return cyberghost.New(storage, randSource)
	case providers.Expressvpn:
		return expressvpn.New(storage, randSource, unzipper, updaterWarner)
	case providers.Fastestvpn:
		return fastestvpn.New(storage, randSource, unzipper, updaterWarner)
	case providers.HideMyAss:
		return hidemyass.New(storage, randSource, client, updaterWarner)
	case providers.Ipvanish:
		return ipvanish.New(storage, randSource, unzipper, updaterWarner)
	case providers.Ivpn:
		return ivpn.New(storage, randSource, client, updaterWarner)
	case providers.Mullvad:
		return mullvad.New(storage, randSource, client)
	case providers.Nordvpn:
		return nordvpn.New(storage, randSource, client, updaterWarner)
	case providers.Perfectprivacy:
		return perfectprivacy.New(storage, randSource, unzipper, updaterWarner)
	case providers.Privado:
		return privado.New(storage, randSource, client, unzipper, updaterWarner)
	case providers.PrivateInternetAccess:
		return privateinternetaccess.New(storage, randSource, timeNow, client)
	case providers.Privatevpn:
		return privatevpn.New(storage, randSource, unzipper, updaterWarner)
	case providers.Protonvpn:
		return protonvpn.New(storage, randSource, client, updaterWarner)
	case providers.Purevpn:
		return purevpn.New(storage, randSource, client, unzipper, updaterWarner)
	case providers.Surfshark:
		return surfshark.New(storage, randSource, client, unzipper, updaterWarner)
	case providers.Torguard:
		return torguard.New(storage, randSource, unzipper, updaterWarner)
	case providers.VPNUnlimited:
		return vpnunlimited.New(storage, randSource, unzipper, updaterWarner)
	case providers.Vyprvpn:
		return vyprvpn.New(storage, randSource, unzipper, updaterWarner)
	case providers.Wevpn:
		return wevpn.New(storage, randSource, updaterWarner)
	case providers.Windscribe:
		return windscribe.New(storage, randSource, client, updaterWarner)
	default:
		panic("provider " + provider + " is unknown") // should never occur
	}
}
