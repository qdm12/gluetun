package storage

import (
	"strconv"
	"time"

	"github.com/qdm12/gluetun/internal/models"
)

func (s *Storage) logVersionDiff(provider string, diff int) {
	diffString := strconv.Itoa(diff)

	message := provider + " servers from file discarded because they are " +
		diffString + " version"
	if diff > 1 {
		message += "s"
	}
	message += " behind"
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
		Version:      hardcoded.Version,
		Cyberghost:   s.mergeCyberghost(hardcoded.Cyberghost, persisted.Cyberghost),
		Fastestvpn:   s.mergeFastestvpn(hardcoded.Fastestvpn, persisted.Fastestvpn),
		HideMyAss:    s.mergeHideMyAss(hardcoded.HideMyAss, persisted.HideMyAss),
		Ipvanish:     s.mergeIpvanish(hardcoded.Ipvanish, persisted.Ipvanish),
		Ivpn:         s.mergeIvpn(hardcoded.Ivpn, persisted.Ivpn),
		Mullvad:      s.mergeMullvad(hardcoded.Mullvad, persisted.Mullvad),
		Nordvpn:      s.mergeNordVPN(hardcoded.Nordvpn, persisted.Nordvpn),
		Privado:      s.mergePrivado(hardcoded.Privado, persisted.Privado),
		Pia:          s.mergePIA(hardcoded.Pia, persisted.Pia),
		Privatevpn:   s.mergePrivatevpn(hardcoded.Privatevpn, persisted.Privatevpn),
		Protonvpn:    s.mergeProtonvpn(hardcoded.Protonvpn, persisted.Protonvpn),
		Purevpn:      s.mergePureVPN(hardcoded.Purevpn, persisted.Purevpn),
		Surfshark:    s.mergeSurfshark(hardcoded.Surfshark, persisted.Surfshark),
		Torguard:     s.mergeTorguard(hardcoded.Torguard, persisted.Torguard),
		VPNUnlimited: s.mergeVPNUnlimited(hardcoded.VPNUnlimited, persisted.VPNUnlimited),
		Vyprvpn:      s.mergeVyprvpn(hardcoded.Vyprvpn, persisted.Vyprvpn),
		Wevpn:        s.mergeWevpn(hardcoded.Wevpn, persisted.Wevpn),
		Windscribe:   s.mergeWindscribe(hardcoded.Windscribe, persisted.Windscribe),
	}
}

func (s *Storage) mergeCyberghost(hardcoded, persisted models.CyberghostServers) models.CyberghostServers {
	if persisted.Timestamp <= hardcoded.Timestamp {
		return hardcoded
	}

	versionDiff := int(hardcoded.Version) - int(persisted.Version)
	if versionDiff > 0 {
		s.logVersionDiff("Cyberghost", versionDiff)
		return hardcoded
	}

	s.logTimeDiff("Cyberghost", persisted.Timestamp, hardcoded.Timestamp)
	return persisted
}

func (s *Storage) mergeFastestvpn(hardcoded, persisted models.FastestvpnServers) models.FastestvpnServers {
	if persisted.Timestamp <= hardcoded.Timestamp {
		return hardcoded
	}
	versionDiff := int(hardcoded.Version) - int(persisted.Version)
	if versionDiff > 0 {
		s.logVersionDiff("FastestVPN", versionDiff)
		return hardcoded
	}
	s.logTimeDiff("FastestVPN", persisted.Timestamp, hardcoded.Timestamp)
	return persisted
}

func (s *Storage) mergeHideMyAss(hardcoded, persisted models.HideMyAssServers) models.HideMyAssServers {
	if persisted.Timestamp <= hardcoded.Timestamp {
		return hardcoded
	}
	versionDiff := int(hardcoded.Version) - int(persisted.Version)
	if versionDiff > 0 {
		s.logVersionDiff("HideMyAss", versionDiff)
		return hardcoded
	}
	s.logTimeDiff("HideMyAss", persisted.Timestamp, hardcoded.Timestamp)
	return persisted
}

func (s *Storage) mergeIpvanish(hardcoded, persisted models.IpvanishServers) models.IpvanishServers {
	if persisted.Timestamp <= hardcoded.Timestamp {
		return hardcoded
	}
	versionDiff := int(hardcoded.Version) - int(persisted.Version)
	if versionDiff > 0 {
		s.logVersionDiff("Ipvanish", versionDiff)
		return hardcoded
	}
	s.logTimeDiff("Ipvanish", persisted.Timestamp, hardcoded.Timestamp)
	return persisted
}

func (s *Storage) mergeIvpn(hardcoded, persisted models.IvpnServers) models.IvpnServers {
	if persisted.Timestamp <= hardcoded.Timestamp {
		return hardcoded
	}
	versionDiff := int(hardcoded.Version) - int(persisted.Version)
	if versionDiff > 0 {
		s.logVersionDiff("Ivpn", versionDiff)
		return hardcoded
	}
	s.logTimeDiff("Ivpn", persisted.Timestamp, hardcoded.Timestamp)
	return persisted
}

func (s *Storage) mergeMullvad(hardcoded, persisted models.MullvadServers) models.MullvadServers {
	if persisted.Timestamp <= hardcoded.Timestamp {
		return hardcoded
	}

	versionDiff := int(hardcoded.Version) - int(persisted.Version)
	if versionDiff > 0 {
		s.logVersionDiff("Mullvad", versionDiff)
		return hardcoded
	}

	s.logTimeDiff("Mullvad", persisted.Timestamp, hardcoded.Timestamp)
	return persisted
}

func (s *Storage) mergeNordVPN(hardcoded, persisted models.NordvpnServers) models.NordvpnServers {
	if persisted.Timestamp <= hardcoded.Timestamp {
		return hardcoded
	}

	versionDiff := int(hardcoded.Version) - int(persisted.Version)
	if versionDiff > 0 {
		s.logVersionDiff("NordVPN", versionDiff)
		return hardcoded
	}

	s.logTimeDiff("NordVPN", persisted.Timestamp, hardcoded.Timestamp)
	return persisted
}

func (s *Storage) mergePrivado(hardcoded, persisted models.PrivadoServers) models.PrivadoServers {
	if persisted.Timestamp <= hardcoded.Timestamp {
		return hardcoded
	}
	versionDiff := int(hardcoded.Version) - int(persisted.Version)
	if versionDiff > 0 {
		s.logVersionDiff("Privado", versionDiff)
		return hardcoded
	}

	s.logTimeDiff("Privado", persisted.Timestamp, hardcoded.Timestamp)
	return persisted
}

func (s *Storage) mergePIA(hardcoded, persisted models.PiaServers) models.PiaServers {
	if persisted.Timestamp <= hardcoded.Timestamp {
		return hardcoded
	}
	versionDiff := int(hardcoded.Version) - int(persisted.Version)
	if versionDiff > 0 {
		s.logVersionDiff("Private Internet Access", versionDiff)
		return hardcoded
	}

	s.logTimeDiff("Private Internet Access", persisted.Timestamp, hardcoded.Timestamp)
	return persisted
}

func (s *Storage) mergePrivatevpn(hardcoded, persisted models.PrivatevpnServers) models.PrivatevpnServers {
	if persisted.Timestamp <= hardcoded.Timestamp {
		return hardcoded
	}
	versionDiff := int(hardcoded.Version) - int(persisted.Version)
	if versionDiff > 0 {
		s.logVersionDiff("PrivateVPN", versionDiff)
		return hardcoded
	}

	s.logTimeDiff("PrivateVPN", persisted.Timestamp, hardcoded.Timestamp)
	return persisted
}

func (s *Storage) mergeProtonvpn(hardcoded, persisted models.ProtonvpnServers) models.ProtonvpnServers {
	if persisted.Timestamp <= hardcoded.Timestamp {
		return hardcoded
	}
	versionDiff := int(hardcoded.Version) - int(persisted.Version)
	if versionDiff > 0 {
		s.logVersionDiff("ProtonVPN", versionDiff)
		return hardcoded
	}

	s.logTimeDiff("ProtonVPN", persisted.Timestamp, hardcoded.Timestamp)
	return persisted
}

func (s *Storage) mergePureVPN(hardcoded, persisted models.PurevpnServers) models.PurevpnServers {
	if persisted.Timestamp <= hardcoded.Timestamp {
		return hardcoded
	}

	versionDiff := int(hardcoded.Version) - int(persisted.Version)
	if versionDiff > 0 {
		s.logVersionDiff("PureVPN", versionDiff)
		return hardcoded
	}

	s.logTimeDiff("PureVPN", persisted.Timestamp, hardcoded.Timestamp)
	return persisted
}

func (s *Storage) mergeSurfshark(hardcoded, persisted models.SurfsharkServers) models.SurfsharkServers {
	if persisted.Timestamp <= hardcoded.Timestamp {
		return hardcoded
	}

	versionDiff := int(hardcoded.Version) - int(persisted.Version)
	if versionDiff > 0 {
		s.logVersionDiff("Surfshark", versionDiff)
		return hardcoded
	}

	s.logTimeDiff("Surfshark", persisted.Timestamp, hardcoded.Timestamp)
	return persisted
}

func (s *Storage) mergeTorguard(hardcoded, persisted models.TorguardServers) models.TorguardServers {
	if persisted.Timestamp <= hardcoded.Timestamp {
		return hardcoded
	}
	versionDiff := int(hardcoded.Version) - int(persisted.Version)
	if versionDiff > 0 {
		s.logVersionDiff("Torguard", versionDiff)
		return hardcoded
	}

	s.logTimeDiff("Torguard", persisted.Timestamp, hardcoded.Timestamp)
	return persisted
}

func (s *Storage) mergeVPNUnlimited(hardcoded, persisted models.VPNUnlimitedServers) models.VPNUnlimitedServers {
	if persisted.Timestamp <= hardcoded.Timestamp {
		return hardcoded
	}
	versionDiff := int(hardcoded.Version) - int(persisted.Version)
	if versionDiff > 0 {
		s.logVersionDiff("VPN Unlimited", versionDiff)
		return hardcoded
	}

	s.logTimeDiff("VPN Unlimited", persisted.Timestamp, hardcoded.Timestamp)
	return persisted
}

func (s *Storage) mergeVyprvpn(hardcoded, persisted models.VyprvpnServers) models.VyprvpnServers {
	if persisted.Timestamp <= hardcoded.Timestamp {
		return hardcoded
	}

	versionDiff := int(hardcoded.Version) - int(persisted.Version)
	if versionDiff > 0 {
		s.logVersionDiff("VyprVPN", versionDiff)
		return hardcoded
	}

	s.logTimeDiff("VyprVPN", persisted.Timestamp, hardcoded.Timestamp)
	return persisted
}

func (s *Storage) mergeWevpn(hardcoded, persisted models.WevpnServers) models.WevpnServers {
	if persisted.Timestamp <= hardcoded.Timestamp {
		return hardcoded
	}

	versionDiff := int(hardcoded.Version) - int(persisted.Version)
	if versionDiff > 0 {
		s.logVersionDiff("WeVPN", versionDiff)
		return hardcoded
	}

	s.logTimeDiff("WeVPN", persisted.Timestamp, hardcoded.Timestamp)
	return persisted
}

func (s *Storage) mergeWindscribe(hardcoded, persisted models.WindscribeServers) models.WindscribeServers {
	if persisted.Timestamp <= hardcoded.Timestamp {
		return hardcoded
	}

	versionDiff := int(hardcoded.Version) - int(persisted.Version)
	if versionDiff > 0 {
		s.logVersionDiff("Windscribe", versionDiff)
		return hardcoded
	}

	s.logTimeDiff("Windscribe", persisted.Timestamp, hardcoded.Timestamp)
	return persisted
}
