package updater

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gluetun/internal/models"
	cyberghost "github.com/qdm12/gluetun/internal/provider/cyberghost/updater"
	expressvpn "github.com/qdm12/gluetun/internal/provider/expressvpn/updater"
	fastestvpn "github.com/qdm12/gluetun/internal/provider/fastestvpn/updater"
	hidemyass "github.com/qdm12/gluetun/internal/provider/hidemyass/updater"
	ipvanish "github.com/qdm12/gluetun/internal/provider/ipvanish/updater"
	ivpn "github.com/qdm12/gluetun/internal/provider/ivpn/updater"
	mullvad "github.com/qdm12/gluetun/internal/provider/mullvad/updater"
	nordvpn "github.com/qdm12/gluetun/internal/provider/nordvpn/updater"
	perfectprivacy "github.com/qdm12/gluetun/internal/provider/perfectprivacy/updater"
	privado "github.com/qdm12/gluetun/internal/provider/privado/updater"
	privateinternetaccess "github.com/qdm12/gluetun/internal/provider/privateinternetaccess/updater"
	privatevpn "github.com/qdm12/gluetun/internal/provider/privatevpn/updater"
	protonvpn "github.com/qdm12/gluetun/internal/provider/protonvpn/updater"
	purevpn "github.com/qdm12/gluetun/internal/provider/purevpn/updater"
	surfshark "github.com/qdm12/gluetun/internal/provider/surfshark/updater"
	torguard "github.com/qdm12/gluetun/internal/provider/torguard/updater"
	vpnunlimited "github.com/qdm12/gluetun/internal/provider/vpnunlimited/updater"
	vyprvpn "github.com/qdm12/gluetun/internal/provider/vyprvpn/updater"
	wevpn "github.com/qdm12/gluetun/internal/provider/wevpn/updater"
	windscribe "github.com/qdm12/gluetun/internal/provider/windscribe/updater"
)

func (u *Updater) updateProvider(ctx context.Context, provider string) (err error) {
	existingServers := u.getProviderServers(provider)
	minServers := getMinServers(existingServers)
	servers, err := u.getServers(ctx, provider, minServers)
	if err != nil {
		return err
	}

	if reflect.DeepEqual(existingServers, servers) {
		return nil
	}

	u.patchProvider(provider, servers)
	return nil
}

func (u *Updater) getServers(ctx context.Context, provider string,
	minServers int) (servers []models.Server, err error) {
	var providerUpdater interface {
		GetServers(ctx context.Context, minServers int) (servers []models.Server, err error)
	}
	switch provider {
	case providers.Custom:
		panic("cannot update custom provider")
	case providers.Cyberghost:
		providerUpdater = cyberghost.New(u.presolver)
	case providers.Expressvpn:
		providerUpdater = expressvpn.New(u.unzipper, u.presolver, u.logger)
	case providers.Fastestvpn:
		providerUpdater = fastestvpn.New(u.unzipper, u.presolver, u.logger)
	case providers.HideMyAss:
		providerUpdater = hidemyass.New(u.client, u.presolver, u.logger)
	case providers.Ipvanish:
		providerUpdater = ipvanish.New(u.unzipper, u.presolver, u.logger)
	case providers.Ivpn:
		providerUpdater = ivpn.New(u.client, u.presolver, u.logger)
	case providers.Mullvad:
		providerUpdater = mullvad.New(u.client)
	case providers.Nordvpn:
		providerUpdater = nordvpn.New(u.client, u.logger)
	case providers.Perfectprivacy:
		providerUpdater = perfectprivacy.New(u.unzipper, u.logger)
	case providers.Privado:
		providerUpdater = privado.New(u.client, u.unzipper, u.presolver, u.logger)
	case providers.PrivateInternetAccess:
		providerUpdater = privateinternetaccess.New(u.client)
	case providers.Privatevpn:
		providerUpdater = privatevpn.New(u.unzipper, u.presolver, u.logger)
	case providers.Protonvpn:
		providerUpdater = protonvpn.New(u.client, u.logger)
	case providers.Purevpn:
		providerUpdater = purevpn.New(u.client, u.unzipper, u.presolver, u.logger)
	case providers.Surfshark:
		providerUpdater = surfshark.New(u.client, u.unzipper, u.presolver, u.logger)
	case providers.Torguard:
		providerUpdater = torguard.New(u.unzipper, u.presolver, u.logger)
	case providers.VPNUnlimited:
		providerUpdater = vpnunlimited.New(u.unzipper, u.presolver, u.logger)
	case providers.Vyprvpn:
		providerUpdater = vyprvpn.New(u.unzipper, u.presolver, u.logger)
	case providers.Wevpn:
		providerUpdater = wevpn.New(u.presolver, u.logger)
	case providers.Windscribe:
		providerUpdater = windscribe.New(u.client, u.logger)
	default:
		panic("provider " + provider + " is unknown")
	}

	servers, err = providerUpdater.GetServers(ctx, minServers)
	return servers, err
}

func (u *Updater) getProviderServers(provider string) (servers []models.Server) {
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

func (u *Updater) patchProvider(provider string, servers []models.Server) {
	providerServers, ok := u.servers.ProviderToServers[provider]
	if !ok {
		panic(fmt.Sprintf("provider %s is unknown", provider))
	}
	providerServers.Timestamp = time.Now().Unix()
	providerServers.Servers = servers
	u.servers.ProviderToServers[provider] = providerServers
}
