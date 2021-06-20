package updater

import (
	"context"
	"fmt"

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
	"github.com/qdm12/gluetun/internal/updater/providers/windscribe"
)

func (u *updater) updateCyberghost(ctx context.Context) (err error) {
	minServers := getMinServers(len(u.servers.Cyberghost.Servers))
	servers, err := cyberghost.GetServers(ctx, u.presolver, minServers)
	if err != nil {
		return err
	}
	if u.options.Stdout {
		u.println(cyberghost.Stringify(servers))
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
	if u.options.Stdout {
		u.println(fastestvpn.Stringify(servers))
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
			u.logger.Warn("HideMyAss: %s", warning)
		}
	}
	if err != nil {
		return err
	}
	if u.options.Stdout {
		u.println(hidemyass.Stringify(servers))
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
			u.logger.Warn("Ipvanish: %s", warning)
		}
	}
	if err != nil {
		return err
	}
	if u.options.Stdout {
		u.println(ipvanish.Stringify(servers))
	}
	u.servers.Ipvanish.Timestamp = u.timeNow().Unix()
	u.servers.Ipvanish.Servers = servers
	return nil
}

func (u *updater) updateIvpn(ctx context.Context) (err error) {
	minServers := getMinServers(len(u.servers.Ivpn.Servers))
	servers, warnings, err := ivpn.GetServers(
		ctx, u.unzipper, u.presolver, minServers)
	if u.options.CLI {
		for _, warning := range warnings {
			u.logger.Warn("Ivpn: %s", warning)
		}
	}
	if err != nil {
		return err
	}
	if u.options.Stdout {
		u.println(ivpn.Stringify(servers))
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
	if u.options.Stdout {
		u.println(mullvad.Stringify(servers))
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
			u.logger.Warn("NordVPN: %s", warning)
		}
	}
	if err != nil {
		return err
	}
	if u.options.Stdout {
		u.println(nordvpn.Stringify(servers))
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
	if u.options.Stdout {
		u.println(pia.Stringify(servers))
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
			u.logger.Warn("Privado: %s", warning)
		}
	}
	if err != nil {
		return err
	}
	if u.options.Stdout {
		u.println(privado.Stringify(servers))
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
			u.logger.Warn("PrivateVPN: %s", warning)
		}
	}
	if err != nil {
		return err
	}
	if u.options.Stdout {
		u.println(privatevpn.Stringify(servers))
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
			u.logger.Warn("ProtonVPN: %s", warning)
		}
	}
	if err != nil {
		return err
	}
	if u.options.Stdout {
		u.println(protonvpn.Stringify(servers))
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
			u.logger.Warn("PureVPN: %s", warning)
		}
	}
	if err != nil {
		return fmt.Errorf("cannot update Purevpn servers: %w", err)
	}
	if u.options.Stdout {
		u.println(purevpn.Stringify(servers))
	}
	u.servers.Purevpn.Timestamp = u.timeNow().Unix()
	u.servers.Purevpn.Servers = servers
	return nil
}

func (u *updater) updateSurfshark(ctx context.Context) (err error) {
	minServers := getMinServers(len(u.servers.Surfshark.Servers))
	servers, warnings, err := surfshark.GetServers(
		ctx, u.unzipper, u.presolver, minServers)
	if u.options.CLI {
		for _, warning := range warnings {
			u.logger.Warn("Surfshark: %s", warning)
		}
	}
	if err != nil {
		return err
	}
	if u.options.Stdout {
		u.println(surfshark.Stringify(servers))
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
			u.logger.Warn("Torguard: %s", warning)
		}
	}
	if err != nil {
		return err
	}
	if u.options.Stdout {
		u.println(torguard.Stringify(servers))
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
	if u.options.Stdout {
		u.println(vpnunlimited.Stringify(servers))
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
			u.logger.Warn("VyprVPN: %s", warning)
		}
	}
	if err != nil {
		return err
	}
	if u.options.Stdout {
		u.println(vyprvpn.Stringify(servers))
	}
	u.servers.Vyprvpn.Timestamp = u.timeNow().Unix()
	u.servers.Vyprvpn.Servers = servers
	return nil
}

func (u *updater) updateWindscribe(ctx context.Context) (err error) {
	minServers := getMinServers(len(u.servers.Windscribe.Servers))
	servers, err := windscribe.GetServers(ctx, u.client, minServers)
	if err != nil {
		return err
	}
	if u.options.Stdout {
		u.println(windscribe.Stringify(servers))
	}
	u.servers.Windscribe.Timestamp = u.timeNow().Unix()
	u.servers.Windscribe.Servers = servers
	return nil
}

func getMinServers(existingServers int) (minServers int) {
	const minRatio = 0.8
	return int(minRatio * float64(existingServers))
}
