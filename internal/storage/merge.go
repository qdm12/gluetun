package storage

import (
	"strconv"
	"time"

	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gluetun/internal/models"
)

func (s *Storage) logVersionDiff(provider string, hardcodedVersion, persistedVersion uint16) {
	message := provider + " servers from file discarded because they have version " +
		strconv.Itoa(int(persistedVersion)) +
		" and hardcoded servers have version " +
		strconv.Itoa(int(hardcodedVersion))
	s.logger.Info(message)
}

func (s *Storage) logTimeDiff(provider string, persistedUnix, hardcodedUnix int64) {
	diff := time.Unix(persistedUnix, 0).Sub(time.Unix(hardcodedUnix, 0))
	if diff < 0 {
		diff = -diff
	}
	diff = diff.Truncate(time.Second)
	message := "Using " + provider + " servers from file which are " +
		diff.String() + " more recent"
	s.logger.Info(message)
}

func (s *Storage) mergeServers(hardcoded, persisted models.AllServers) models.AllServers {
	return models.AllServers{
		Version:        hardcoded.Version,
		Cyberghost:     s.mergeProviderServers(providers.Cyberghost, hardcoded.Cyberghost, persisted.Cyberghost),
		Expressvpn:     s.mergeProviderServers(providers.Expressvpn, hardcoded.Expressvpn, persisted.Expressvpn),
		Fastestvpn:     s.mergeProviderServers(providers.Fastestvpn, hardcoded.Fastestvpn, persisted.Fastestvpn),
		HideMyAss:      s.mergeProviderServers(providers.HideMyAss, hardcoded.HideMyAss, persisted.HideMyAss),
		Ipvanish:       s.mergeProviderServers(providers.Ipvanish, hardcoded.Ipvanish, persisted.Ipvanish),
		Ivpn:           s.mergeProviderServers(providers.Ivpn, hardcoded.Ivpn, persisted.Ivpn),
		Mullvad:        s.mergeProviderServers(providers.Mullvad, hardcoded.Mullvad, persisted.Mullvad),
		Nordvpn:        s.mergeProviderServers(providers.Nordvpn, hardcoded.Nordvpn, persisted.Nordvpn),
		Perfectprivacy: s.mergeProviderServers(providers.Perfectprivacy, hardcoded.Perfectprivacy, persisted.Perfectprivacy),
		Privado:        s.mergeProviderServers(providers.Privado, hardcoded.Privado, persisted.Privado),
		Pia:            s.mergeProviderServers(providers.PrivateInternetAccess, hardcoded.Pia, persisted.Pia),
		Privatevpn:     s.mergeProviderServers(providers.Privatevpn, hardcoded.Privatevpn, persisted.Privatevpn),
		Protonvpn:      s.mergeProviderServers(providers.Protonvpn, hardcoded.Protonvpn, persisted.Protonvpn),
		Purevpn:        s.mergeProviderServers(providers.Purevpn, hardcoded.Purevpn, persisted.Purevpn),
		Surfshark:      s.mergeProviderServers(providers.Surfshark, hardcoded.Surfshark, persisted.Surfshark),
		Torguard:       s.mergeProviderServers(providers.Torguard, hardcoded.Torguard, persisted.Torguard),
		VPNUnlimited:   s.mergeProviderServers(providers.VPNUnlimited, hardcoded.VPNUnlimited, persisted.VPNUnlimited),
		Vyprvpn:        s.mergeProviderServers(providers.Vyprvpn, hardcoded.Vyprvpn, persisted.Vyprvpn),
		Wevpn:          s.mergeProviderServers(providers.Wevpn, hardcoded.Wevpn, persisted.Wevpn),
		Windscribe:     s.mergeProviderServers(providers.Windscribe, hardcoded.Windscribe, persisted.Windscribe),
	}
}

func (s *Storage) mergeProviderServers(provider string,
	hardcoded, persisted models.Servers) (merged models.Servers) {
	if persisted.Timestamp <= hardcoded.Timestamp {
		return hardcoded
	}

	s.logTimeDiff(provider, persisted.Timestamp, hardcoded.Timestamp)
	return persisted
}
