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
		Pia: models.PiaServers{
			Version:   3,
			Timestamp: 1613479739,
			Servers:   PIAServers(),
		},
		Purevpn: models.PurevpnServers{
			Version:   1,
			Timestamp: 1612031135,
			Servers:   PurevpnServers(),
		},
		Privado: models.PrivadoServers{
			Version:   2,
			Timestamp: 1612031135,
			Servers:   PrivadoServers(),
		},
		Surfshark: models.SurfsharkServers{
			Version:   1,
			Timestamp: 1612031135,
			Servers:   SurfsharkServers(),
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
