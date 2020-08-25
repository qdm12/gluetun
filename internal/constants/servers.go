package constants

import "github.com/qdm12/gluetun/internal/models"

func GetAllServers() (allServers models.AllServers) {
	return models.AllServers{
		Version: 1, // used for migration of the top level scheme
		Cyberghost: models.CyberghostServers{
			Version:   1,          // model version
			Timestamp: 1598236838, // latest takes precedence
			Servers:   CyberghostServers(),
		},
		Mullvad: models.MullvadServers{
			Version:   1,
			Timestamp: 1598236838,
			Servers:   MullvadServers(),
		},
		Nordvpn: models.NordvpnServers{
			Version:   1,
			Timestamp: 1598236838,
			Servers:   NordvpnServers(),
		},
		Pia: models.PiaServers{
			Version:   1,
			Timestamp: 1598236838,
			Servers:   PIAServers(),
		},
		Purevpn: models.PurevpnServers{
			Version:   1,
			Timestamp: 1598236838,
			Servers:   PurevpnServers(),
		},
		Surfshark: models.SurfsharkServers{
			Version:   1,
			Timestamp: 1598236838,
			Servers:   SurfsharkServers(),
		},
		Vyprvpn: models.VyprvpnServers{
			Version:   1,
			Timestamp: 1598236838,
			Servers:   VyprvpnServers(),
		},
		Windscribe: models.WindscribeServers{
			Version:   1,
			Timestamp: 1598236838,
			Servers:   WindscribeServers(),
		},
	}
}
