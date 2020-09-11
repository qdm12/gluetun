package constants

import "github.com/qdm12/gluetun/internal/models"

func GetAllServers() (allServers models.AllServers) {
	return models.AllServers{
		Version: 1, // used for migration of the top level scheme
		Cyberghost: models.CyberghostServers{
			Version:   1,          // model version
			Timestamp: 1599323261, // latest takes precedence
			Servers:   CyberghostServers(),
		},
		Mullvad: models.MullvadServers{
			Version:   1,
			Timestamp: 1599323261,
			Servers:   MullvadServers(),
		},
		Nordvpn: models.NordvpnServers{
			Version:   1,
			Timestamp: 1599323261,
			Servers:   NordvpnServers(),
		},
		Pia: models.PiaServers{
			Version:   1,
			Timestamp: 1599323261,
			Servers:   PIAServers(),
		},
		PiaOld: models.PiaServers{
			Version:   1,
			Timestamp: 1599786395,
			Servers:   PIAOldServers(),
		},
		Purevpn: models.PurevpnServers{
			Version:   1,
			Timestamp: 1599323261,
			Servers:   PurevpnServers(),
		},
		Surfshark: models.SurfsharkServers{
			Version:   1,
			Timestamp: 1599323261,
			Servers:   SurfsharkServers(),
		},
		Vyprvpn: models.VyprvpnServers{
			Version:   1,
			Timestamp: 1599323261,
			Servers:   VyprvpnServers(),
		},
		Windscribe: models.WindscribeServers{
			Version:   1,
			Timestamp: 1599323261,
			Servers:   WindscribeServers(),
		},
	}
}
