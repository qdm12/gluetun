package updater

import (
	"context"
	"fmt"
	"reflect"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/updater/providers/cyberghost"
	"github.com/qdm12/gluetun/internal/updater/providers/fastestvpn"
	"github.com/qdm12/gluetun/internal/updater/providers/hidemyass"
	"github.com/qdm12/gluetun/internal/updater/providers/ipvanish"
	"github.com/qdm12/gluetun/internal/updater/providers/ivpn"
	"github.com/qdm12/gluetun/internal/updater/providers/mullvad"
	"github.com/qdm12/gluetun/internal/updater/providers/nordvpn"
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

func (u *updater) updateCyberghost(ctx context.Context) (err error) {
	minServers := getMinServers(len(u.servers.Cyberghost.Servers))
	servers, err := cyberghost.GetServers(ctx, u.presolver, minServers)
	if err != nil {
		return err
	}

	if reflect.DeepEqual(u.servers.Cyberghost.Servers, servers) {
		return nil
	}

	u.servers.Cyberghost.Timestamp = u.timeNow().Unix()
	u.servers.Cyberghost.Servers = servers
	return nil
}

func (u *updater) updateFastestvpn(ctx context.Context) (err error) {
	minServers := getMinServers(len(u.servers.Fastestvpn.Servers))
	servers, warnings, err := fastestvpn.GetServers(
		ctx, u.unzipper, u.presolver, minServers)
	if u.options.CLI {
		for _, warning := range warnings {
			u.logger.Warn("FastestVPN: " + warning)
		}
	}
	if err != nil {
		return err
	}

	if reflect.DeepEqual(u.servers.Fastestvpn.Servers, servers) {
		return nil
	}

	u.servers.Fastestvpn.Timestamp = u.timeNow().Unix()
	u.servers.Fastestvpn.Servers = servers
	return nil
}

func (u *updater) updateHideMyAss(ctx context.Context) (err error) {
	minServers := getMinServers(len(u.servers.HideMyAss.Servers))
	servers, warnings, err := hidemyass.GetServers(
		ctx, u.client, u.presolver, minServers)
	if u.options.CLI {
		for _, warning := range warnings {
			u.logger.Warn("HideMyAss: " + warning)
		}
	}
	if err != nil {
		return err
	}

	if reflect.DeepEqual(u.servers.HideMyAss.Servers, servers) {
		return nil
	}

	u.servers.HideMyAss.Timestamp = u.timeNow().Unix()
	u.servers.HideMyAss.Servers = servers
	return nil
}

func (u *updater) updateIpvanish(ctx context.Context) (err error) {
	minServers := getMinServers(len(u.servers.Ipvanish.Servers))
	servers, warnings, err := ipvanish.GetServers(
		ctx, u.unzipper, u.presolver, minServers)
	if u.options.CLI {
		for _, warning := range warnings {
			u.logger.Warn("Ipvanish: " + warning)
		}
	}
	if err != nil {
		return err
	}

	if reflect.DeepEqual(u.servers.Ipvanish.Servers, servers) {
		return nil
	}

	u.servers.Ipvanish.Timestamp = u.timeNow().Unix()
	u.servers.Ipvanish.Servers = servers
	return nil
}

func (u *updater) updateIvpn(ctx context.Context) (err error) {
	minServers := getMinServers(len(u.servers.Ivpn.Servers))
	servers, warnings, err := ivpn.GetServers(
		ctx, u.client, u.presolver, minServers)
	if u.options.CLI {
		for _, warning := range warnings {
			u.logger.Warn("Ivpn: " + warning)
		}
	}
	if err != nil {
		return err
	}

	if reflect.DeepEqual(u.servers.Ivpn.Servers, servers) {
		return nil
	}

	u.servers.Ivpn.Timestamp = u.timeNow().Unix()
	u.servers.Ivpn.Servers = servers
	return nil
}

func (u *updater) updateMullvad(ctx context.Context) (err error) {
	minServers := getMinServers(len(u.servers.Mullvad.Servers))
	servers, err := mullvad.GetServers(ctx, u.client, minServers)
	if err != nil {
		return err
	}

	if reflect.DeepEqual(u.servers.Mullvad.Servers, servers) {
		return nil
	}

	u.servers.Mullvad.Timestamp = u.timeNow().Unix()
	u.servers.Mullvad.Servers = servers
	return nil
}

func (u *updater) updateNordvpn(ctx context.Context) (err error) {
	minServers := getMinServers(len(u.servers.Nordvpn.Servers))
	servers, warnings, err := nordvpn.GetServers(ctx, u.client, minServers)
	if u.options.CLI {
		for _, warning := range warnings {
			u.logger.Warn("NordVPN: " + warning)
		}
	}
	if err != nil {
		return err
	}

	if reflect.DeepEqual(u.servers.Nordvpn.Servers, servers) {
		return nil
	}

	u.servers.Nordvpn.Timestamp = u.timeNow().Unix()
	u.servers.Nordvpn.Servers = servers
	return nil
}

func (u *updater) updatePIA(ctx context.Context) (err error) {
	minServers := getMinServers(len(u.servers.Pia.Servers))
	servers, err := pia.GetServers(ctx, u.client, minServers)
	if err != nil {
		return err
	}

	if reflect.DeepEqual(u.servers.Pia.Servers, servers) {
		return nil
	}

	u.servers.Pia.Timestamp = u.timeNow().Unix()
	u.servers.Pia.Servers = servers
	return nil
}

func (u *updater) updatePrivado(ctx context.Context) (err error) {
	minServers := getMinServers(len(u.servers.Privado.Servers))
	servers, warnings, err := privado.GetServers(
		ctx, u.unzipper, u.client, u.presolver, minServers)
	if u.options.CLI {
		for _, warning := range warnings {
			u.logger.Warn("Privado: " + warning)
		}
	}
	if err != nil {
		return err
	}

	if reflect.DeepEqual(u.servers.Privado.Servers, servers) {
		return nil
	}

	u.servers.Privado.Timestamp = u.timeNow().Unix()
	u.servers.Privado.Servers = servers
	return nil
}

func (u *updater) updatePrivatevpn(ctx context.Context) (err error) {
	minServers := getMinServers(len(u.servers.Privatevpn.Servers))
	servers, warnings, err := privatevpn.GetServers(
		ctx, u.unzipper, u.presolver, minServers)
	if u.options.CLI {
		for _, warning := range warnings {
			u.logger.Warn("PrivateVPN: " + warning)
		}
	}
	if err != nil {
		return err
	}

	if reflect.DeepEqual(u.servers.Privatevpn.Servers, servers) {
		return nil
	}

	u.servers.Privatevpn.Timestamp = u.timeNow().Unix()
	u.servers.Privatevpn.Servers = servers
	return nil
}

func (u *updater) updateProtonvpn(ctx context.Context) (err error) {
	minServers := getMinServers(len(u.servers.Privatevpn.Servers))
	servers, warnings, err := protonvpn.GetServers(ctx, u.client, minServers)
	if u.options.CLI {
		for _, warning := range warnings {
			u.logger.Warn("ProtonVPN: " + warning)
		}
	}
	if err != nil {
		return err
	}

	if reflect.DeepEqual(u.servers.Protonvpn.Servers, servers) {
		return nil
	}

	u.servers.Protonvpn.Timestamp = u.timeNow().Unix()
	u.servers.Protonvpn.Servers = servers
	return nil
}

func (u *updater) updatePurevpn(ctx context.Context) (err error) {
	minServers := getMinServers(len(u.servers.Purevpn.Servers))
	servers, warnings, err := purevpn.GetServers(
		ctx, u.client, u.unzipper, u.presolver, minServers)
	if u.options.CLI {
		for _, warning := range warnings {
			u.logger.Warn("PureVPN: " + warning)
		}
	}
	if err != nil {
		return fmt.Errorf("cannot update Purevpn servers: %w", err)
	}

	if reflect.DeepEqual(u.servers.Purevpn.Servers, servers) {
		return nil
	}

	u.servers.Purevpn.Timestamp = u.timeNow().Unix()
	u.servers.Purevpn.Servers = servers
	return nil
}

func (u *updater) updateSurfshark(ctx context.Context) (err error) {
	minServers := getMinServers(len(u.servers.Surfshark.Servers))
	servers, warnings, err := surfshark.GetServers(
		ctx, u.unzipper, u.client, u.presolver, minServers)
	if u.options.CLI {
		for _, warning := range warnings {
			u.logger.Warn("Surfshark: " + warning)
		}
	}
	if err != nil {
		return err
	}

	if reflect.DeepEqual(u.servers.Surfshark.Servers, servers) {
		return nil
	}

	u.servers.Surfshark.Timestamp = u.timeNow().Unix()
	u.servers.Surfshark.Servers = servers
	return nil
}

func (u *updater) updateTorguard(ctx context.Context) (err error) {
	minServers := getMinServers(len(u.servers.Torguard.Servers))
	servers, warnings, err := torguard.GetServers(
		ctx, u.unzipper, u.presolver, minServers)
	if u.options.CLI {
		for _, warning := range warnings {
			u.logger.Warn("Torguard: " + warning)
		}
	}
	if err != nil {
		return err
	}

	if reflect.DeepEqual(u.servers.Torguard.Servers, servers) {
		return nil
	}

	u.servers.Torguard.Timestamp = u.timeNow().Unix()
	u.servers.Torguard.Servers = servers
	return nil
}

func (u *updater) updateVPNUnlimited(ctx context.Context) (err error) {
	minServers := getMinServers(len(u.servers.VPNUnlimited.Servers))
	servers, warnings, err := vpnunlimited.GetServers(
		ctx, u.unzipper, u.presolver, minServers)
	if u.options.CLI {
		for _, warning := range warnings {
			u.logger.Warn(constants.VPNUnlimited + ": " + warning)
		}
	}
	if err != nil {
		return err
	}

	if reflect.DeepEqual(u.servers.VPNUnlimited.Servers, servers) {
		return nil
	}

	u.servers.VPNUnlimited.Timestamp = u.timeNow().Unix()
	u.servers.VPNUnlimited.Servers = servers
	return nil
}

func (u *updater) updateVyprvpn(ctx context.Context) (err error) {
	minServers := getMinServers(len(u.servers.Vyprvpn.Servers))
	servers, warnings, err := vyprvpn.GetServers(
		ctx, u.unzipper, u.presolver, minServers)
	if u.options.CLI {
		for _, warning := range warnings {
			u.logger.Warn("VyprVPN: " + warning)
		}
	}
	if err != nil {
		return err
	}

	if reflect.DeepEqual(u.servers.Vyprvpn.Servers, servers) {
		return nil
	}

	u.servers.Vyprvpn.Timestamp = u.timeNow().Unix()
	u.servers.Vyprvpn.Servers = servers
	return nil
}

func (u *updater) updateWevpn(ctx context.Context) (err error) {
	minServers := getMinServers(len(u.servers.Wevpn.Servers))
	servers, warnings, err := wevpn.GetServers(ctx, u.presolver, minServers)
	if u.options.CLI {
		for _, warning := range warnings {
			u.logger.Warn("WeVPN: " + warning)
		}
	}
	if err != nil {
		return err
	}

	if reflect.DeepEqual(u.servers.Wevpn.Servers, servers) {
		return nil
	}

	u.servers.Wevpn.Timestamp = u.timeNow().Unix()
	u.servers.Wevpn.Servers = servers
	return nil
}

func (u *updater) updateWindscribe(ctx context.Context) (err error) {
	minServers := getMinServers(len(u.servers.Windscribe.Servers))
	servers, err := windscribe.GetServers(ctx, u.client, minServers)
	if err != nil {
		return err
	}

	if reflect.DeepEqual(u.servers.Windscribe.Servers, servers) {
		return nil
	}

	u.servers.Windscribe.Timestamp = u.timeNow().Unix()
	u.servers.Windscribe.Servers = servers
	return nil
}

func getMinServers(existingServers int) (minServers int) {
	const minRatio = 0.8
	return int(minRatio * float64(existingServers))
}
