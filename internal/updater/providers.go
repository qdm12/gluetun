package updater

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/updater/providers/cyberghost"
	"github.com/qdm12/gluetun/internal/updater/providers/expressvpn"
	"github.com/qdm12/gluetun/internal/updater/providers/fastestvpn"
	"github.com/qdm12/gluetun/internal/updater/providers/hidemyass"
	"github.com/qdm12/gluetun/internal/updater/providers/ipvanish"
	"github.com/qdm12/gluetun/internal/updater/providers/ivpn"
	"github.com/qdm12/gluetun/internal/updater/providers/mullvad"
	"github.com/qdm12/gluetun/internal/updater/providers/nordvpn"
	"github.com/qdm12/gluetun/internal/updater/providers/perfectprivacy"
	"github.com/qdm12/gluetun/internal/updater/providers/pia"
	"github.com/qdm12/gluetun/internal/updater/providers/privado"
	"github.com/qdm12/gluetun/internal/updater/providers/privatevpn"
	"github.com/qdm12/gluetun/internal/updater/providers/protonvpn"
	"github.com/qdm12/gluetun/internal/updater/providers/purevpn"
	"github.com/qdm12/gluetun/internal/updater/providers/surfshark"
	"github.com/qdm12/gluetun/internal/updater/providers/torguard"
	"github.com/qdm12/gluetun/internal/updater/providers/vpnunlimited"
	"github.com/qdm12/gluetun/internal/updater/providers/vyprvpn"
	"github.com/qdm12/gluetun/internal/updater/providers/wevpn"
	"github.com/qdm12/gluetun/internal/updater/providers/windscribe"
)

func (u *updater) updateProvider(ctx context.Context, provider string) (
	warnings []string, err error) {
	existingServers := u.getProviderServers(provider)
	minServers := getMinServers(existingServers)
	servers, warnings, err := u.getServers(ctx, provider, minServers)
	if err != nil {
		return warnings, err
	}

	if reflect.DeepEqual(existingServers, servers) {
		return warnings, nil
	}

	u.patchProvider(provider, servers)
	return warnings, nil
}

func (u *updater) getServers(ctx context.Context, provider string,
	minServers int) (servers []models.Server, warnings []string, err error) {
	switch provider {
	case providers.Custom:
		panic("cannot update custom provider")
	case providers.Cyberghost:
		servers, err = cyberghost.GetServers(ctx, u.presolver, minServers)
		return servers, nil, err
	case providers.Expressvpn:
		return expressvpn.GetServers(ctx, u.unzipper, u.presolver, minServers)
	case providers.Fastestvpn:
		return fastestvpn.GetServers(ctx, u.unzipper, u.presolver, minServers)
	case providers.HideMyAss:
		return hidemyass.GetServers(ctx, u.client, u.presolver, minServers)
	case providers.Ipvanish:
		return ipvanish.GetServers(ctx, u.unzipper, u.presolver, minServers)
	case providers.Ivpn:
		return ivpn.GetServers(ctx, u.client, u.presolver, minServers)
	case providers.Mullvad:
		servers, err = mullvad.GetServers(ctx, u.client, minServers)
		return servers, nil, err
	case providers.Nordvpn:
		return nordvpn.GetServers(ctx, u.client, minServers)
	case providers.Perfectprivacy:
		return perfectprivacy.GetServers(ctx, u.unzipper, minServers)
	case providers.Privado:
		return privado.GetServers(ctx, u.unzipper, u.client, u.presolver, minServers)
	case providers.PrivateInternetAccess:
		servers, err = pia.GetServers(ctx, u.client, minServers)
		return servers, nil, err
	case providers.Privatevpn:
		return privatevpn.GetServers(ctx, u.unzipper, u.presolver, minServers)
	case providers.Protonvpn:
		return protonvpn.GetServers(ctx, u.client, minServers)
	case providers.Purevpn:
		return purevpn.GetServers(ctx, u.client, u.unzipper, u.presolver, minServers)
	case providers.Surfshark:
		return surfshark.GetServers(ctx, u.unzipper, u.client, u.presolver, minServers)
	case providers.Torguard:
		return torguard.GetServers(ctx, u.unzipper, u.presolver, minServers)
	case providers.VPNUnlimited:
		return vpnunlimited.GetServers(ctx, u.unzipper, u.presolver, minServers)
	case providers.Vyprvpn:
		return vyprvpn.GetServers(ctx, u.unzipper, u.presolver, minServers)
	case providers.Wevpn:
		return wevpn.GetServers(ctx, u.presolver, minServers)
	case providers.Windscribe:
		servers, err = windscribe.GetServers(ctx, u.client, minServers)
		return servers, nil, err
	default:
		panic("provider " + provider + " is unknown")
	}
}

func (u *updater) getProviderServers(provider string) (servers []models.Server) {
	providerServers, ok := u.servers.ProviderToServers[provider]
	if !ok {
		panic(fmt.Sprintf("provider %s is unknown", provider))
	}
	return providerServers.Servers
}

func getMinServers(servers []models.Server) (minServers int) {
	const minRatio = 0.8
	return int(minRatio * float64(len(servers)))
}

func (u *updater) patchProvider(provider string, servers []models.Server) {
	providerServers, ok := u.servers.ProviderToServers[provider]
	if !ok {
		panic(fmt.Sprintf("provider %s is unknown", provider))
	}
	providerServers.Timestamp = time.Now().Unix()
	providerServers.Servers = servers
	u.servers.ProviderToServers[provider] = providerServers
}
