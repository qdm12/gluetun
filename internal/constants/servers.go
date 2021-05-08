package constants

import "github.com/qdm12/gluetun/internal/models"

func GetAllServers() (allServers models.AllServers) {
	//nolint:gomnd
	return models.AllServers{
		Version: 1, // used for migration of the top level scheme
		Cyberghost: models.CyberghostServers{
			Version:   2,          // model version
			Timestamp: 1620491290, // latest takes precedence
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
		Mullvad: models.MullvadServers{
			Version:   2,
			Timestamp: 1620500848,
			Servers:   MullvadServers(),
		},
		Nordvpn: models.NordvpnServers{
			Version:   1,
			Timestamp: 1620435633,
			Servers:   NordvpnServers(),
		},
		Privado: models.PrivadoServers{
			Version:   2,
			Timestamp: 1620435633,
			Servers:   PrivadoServers(),
		},
		Privatevpn: models.PrivatevpnServers{
			Version:   1,
			Timestamp: 1620435633,
			Servers:   PrivatevpnServers(),
		},
		Protonvpn: models.ProtonvpnServers{
			Version:   1,
			Timestamp: 1620435633,
			Servers:   ProtonvpnServers(),
		},
		Pia: models.PiaServers{
			Version:   4,
			Timestamp: 1620435633,
			Servers:   PIAServers(),
		},
		Purevpn: models.PurevpnServers{
			Version:   1,
			Timestamp: 1620435633,
			Servers:   PurevpnServers(),
		},
		Surfshark: models.SurfsharkServers{
			Version:   1,
			Timestamp: 1620435633,
			Servers:   SurfsharkServers(),
		},
		Torguard: models.TorguardServers{
			Version:   1,
			Timestamp: 1613357861,
			Servers:   TorguardServers(),
		},
		Vyprvpn: models.VyprvpnServers{
			Version:   1,
			Timestamp: 1620326951,
			Servers:   VyprvpnServers(),
		},
		Windscribe: models.WindscribeServers{
			Version:   2,
			Timestamp: 1620435633,
			Servers:   WindscribeServers(),
		},
	}
}
