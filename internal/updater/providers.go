package updater

import (
	"context"
	"fmt"

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
	existingServersCount := u.storage.GetServersCount(provider)
	minServers := getMinServers(existingServersCount)
	servers, err := u.fetchServers(ctx, provider, minServers)
	if err != nil {
		return fmt.Errorf("cannot get servers: %w", err)
	}

	if u.storage.ServersAreEqual(provider, servers) {
		return nil
	}

	// Note the servers variable must NOT BE MUTATED after this call,
	// since the implementation does not deep copy the servers.
	// TODO set in storage in provider updater directly, server by server,
	// to avoid accumulating server data in memory.
	err = u.storage.SetServers(provider, servers)
	if err != nil {
		return fmt.Errorf("cannot set servers to storage: %w", err)
	}
	return nil
}

func (u *Updater) fetchServers(ctx context.Context, provider string,
	minServers int) (servers []models.Server, err error) {
	var providerUpdater interface {
		FetchServers(ctx context.Context, minServers int) (servers []models.Server, err error)
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

	return providerUpdater.FetchServers(ctx, minServers)
}

func getMinServers(existingServersCount int) (minServers int) {
	const minRatio = 0.8
	return int(minRatio * float64(existingServersCount))
}
