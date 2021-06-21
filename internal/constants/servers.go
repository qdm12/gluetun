package constants

import "github.com/qdm12/gluetun/internal/models"

func GetAllServers() (allServers models.AllServers) {
	//nolint:gomnd
	return models.AllServers{
		Version: 1, // used for migration of the top level scheme
		Cyberghost: models.CyberghostServers{
			Version:   2,          // model version
			Timestamp: 1624307338, // latest takes precedence
			Servers:   CyberghostServers(),
		},
		Fastestvpn: models.FastestvpnServers{
			Version:   1,
			Timestamp: 1620435633,
			Servers:   FastestvpnServers(),
		},
		HideMyAss: models.HideMyAssServers{
			Version:   1,
			Timestamp: 1620435633,
			Servers:   HideMyAssServers(),
		},
		Ipvanish: models.IpvanishServers{
			Version:   1,
			Timestamp: 1622430497,
			Servers:   IpvanishServers(),
		},
		Ivpn: models.IvpnServers{
			Version:   1,
			Timestamp: 1624120443,
			Servers:   IvpnServers(),
		},
		Mullvad: models.MullvadServers{
			Version:   2,
			Timestamp: 1620500848,
			Servers:   MullvadServers(),
		},
		Nordvpn: models.NordvpnServers{
			Version:   2,
			Timestamp: 1620514180,
			Servers:   NordvpnServers(),
		},
		Privado: models.PrivadoServers{
			Version:   3,
			Timestamp: 1620520278,
			Servers:   PrivadoServers(),
		},
		Privatevpn: models.PrivatevpnServers{
			Version:   1,
			Timestamp: 1620435633,
			Servers:   PrivatevpnServers(),
		},
		Protonvpn: models.ProtonvpnServers{
			Version:   1,
			Timestamp: 1621791438,
			Servers:   ProtonvpnServers(),
		},
		Pia: models.PiaServers{
			Version:   6,
			Timestamp: 1620663401,
			Servers:   PIAServers(),
		},
		Purevpn: models.PurevpnServers{
			Version:   2,
			Timestamp: 1622644308,
			Servers:   PurevpnServers(),
		},
		Surfshark: models.SurfsharkServers{
			Version:   2,
			Timestamp: 1620607876,
			Servers:   SurfsharkServers(),
		},
		Torguard: models.TorguardServers{
			Version:   2,
			Timestamp: 1620611129,
			Servers:   TorguardServers(),
		},
		VPNUnlimited: models.VPNUnlimitedServers{
			Version:   1,
			Timestamp: 1623950304,
			Servers:   VPNUnlimitedServers(),
		},
		Vyprvpn: models.VyprvpnServers{
			Version:   2,
			Timestamp: 1620612506,
			Servers:   VyprvpnServers(),
		},
		Windscribe: models.WindscribeServers{
			Version:   3,
			Timestamp: 1620657134,
			Servers:   WindscribeServers(),
		},
	}
}
