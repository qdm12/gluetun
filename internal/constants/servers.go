package constants

import "github.com/qdm12/gluetun/internal/models"

func GetAllServers() (allServers models.AllServers) {
	//nolint:gomnd
	return models.AllServers{
		Version: 1, // used for migration of the top level scheme
		Cyberghost: models.CyberghostServers{
			Version:   1,          // model version
			Timestamp: 1612031135, // latest takes precedence
			Servers:   CyberghostServers(),
		},
		HideMyAss: models.HideMyAssServers{
			Version:   1,
			Timestamp: 1614562368,
			Servers:   HideMyAssServers(),
		},
		Mullvad: models.MullvadServers{
			Version:   1,
			Timestamp: 1612031135,
			Servers:   MullvadServers(),
		},
		Nordvpn: models.NordvpnServers{
			Version:   1,
			Timestamp: 1611096594,
			Servers:   NordvpnServers(),
		},
		Privado: models.PrivadoServers{
			Version:   2,
			Timestamp: 1612031135,
			Servers:   PrivadoServers(),
		},
		Pia: models.PiaServers{
			Version:   4,
			Timestamp: 1613480675,
			Servers:   PIAServers(),
		},
		Purevpn: models.PurevpnServers{
			Version:   1,
			Timestamp: 1612031135,
			Servers:   PurevpnServers(),
		},
		Surfshark: models.SurfsharkServers{
			Version:   1,
			Timestamp: 1612031135,
			Servers:   SurfsharkServers(),
		},
		Torguard: models.TorguardServers{
			Version:   1,
			Timestamp: 1613357861,
			Servers:   TorguardServers(),
		},
		Vyprvpn: models.VyprvpnServers{
			Version:   1,
			Timestamp: 1612031135,
			Servers:   VyprvpnServers(),
		},
		Windscribe: models.WindscribeServers{
			Version:   2,
			Timestamp: 1612031135,
			Servers:   WindscribeServers(),
		},
	}
}
