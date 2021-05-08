package storage

import (
	"strconv"
	"time"

	"github.com/qdm12/gluetun/internal/models"
)

func getUnixTimeDifference(unix1, unix2 int64) (difference time.Duration) {
	difference = time.Unix(unix1, 0).Sub(time.Unix(unix2, 0))
	if difference < 0 {
		difference = -difference
	}
	return difference.Truncate(time.Second)
}

func (s *storage) logVersionDiff(provider string, diff uint16) {
	diffString := strconv.Itoa(int(diff))

	message := provider + " servers from file discarded because they are " +
		diffString + " version"
	if diff > 1 {
		message += "s"
	}
	s.logger.Info(message)
}

func (s *storage) mergeServers(hardcoded, persisted models.AllServers) models.AllServers {
	return models.AllServers{
		Version:    hardcoded.Version,
		Cyberghost: s.mergeCyberghost(hardcoded.Cyberghost, persisted.Cyberghost),
		Fastestvpn: s.mergeFastestvpn(hardcoded.Fastestvpn, persisted.Fastestvpn),
		HideMyAss:  s.mergeHideMyAss(hardcoded.HideMyAss, persisted.HideMyAss),
		Mullvad:    s.mergeMullvad(hardcoded.Mullvad, persisted.Mullvad),
		Nordvpn:    s.mergeNordVPN(hardcoded.Nordvpn, persisted.Nordvpn),
		Privado:    s.mergePrivado(hardcoded.Privado, persisted.Privado),
		Pia:        s.mergePIA(hardcoded.Pia, persisted.Pia),
		Privatevpn: s.mergePrivatevpn(hardcoded.Privatevpn, persisted.Privatevpn),
		Protonvpn:  s.mergeProtonvpn(hardcoded.Protonvpn, persisted.Protonvpn),
		Purevpn:    s.mergePureVPN(hardcoded.Purevpn, persisted.Purevpn),
		Surfshark:  s.mergeSurfshark(hardcoded.Surfshark, persisted.Surfshark),
		Torguard:   s.mergeTorguard(hardcoded.Torguard, persisted.Torguard),
		Vyprvpn:    s.mergeVyprvpn(hardcoded.Vyprvpn, persisted.Vyprvpn),
		Windscribe: s.mergeWindscribe(hardcoded.Windscribe, persisted.Windscribe),
	}
}

func (s *storage) mergeCyberghost(hardcoded, persisted models.CyberghostServers) models.CyberghostServers {
	if persisted.Timestamp <= hardcoded.Timestamp {
		return hardcoded
	}

	versionDiff := hardcoded.Version - persisted.Version
	if versionDiff > 0 {
		s.logVersionDiff("Cyberghost", versionDiff)
		return hardcoded
	}

	s.logger.Info("Using Cyberghost servers from file (%s more recent)",
		getUnixTimeDifference(persisted.Timestamp, hardcoded.Timestamp))
	return persisted
}

func (s *storage) mergeFastestvpn(hardcoded, persisted models.FastestvpnServers) models.FastestvpnServers {
	if persisted.Timestamp <= hardcoded.Timestamp {
		return hardcoded
	}
	versionDiff := hardcoded.Version - persisted.Version
	if versionDiff > 0 {
		s.logVersionDiff("FastestVPN", versionDiff)
		return hardcoded
	}
	s.logger.Info("Using Fastestvpn servers from file (%s more recent)",
		getUnixTimeDifference(persisted.Timestamp, hardcoded.Timestamp))
	return persisted
}

func (s *storage) mergeHideMyAss(hardcoded, persisted models.HideMyAssServers) models.HideMyAssServers {
	if persisted.Timestamp <= hardcoded.Timestamp {
		return hardcoded
	}
	versionDiff := hardcoded.Version - persisted.Version
	if versionDiff > 0 {
		s.logVersionDiff("HideMyAss", versionDiff)
		return hardcoded
	}
	s.logger.Info("Using HideMyAss servers from file (%s more recent)",
		getUnixTimeDifference(persisted.Timestamp, hardcoded.Timestamp))
	return persisted
}

func (s *storage) mergeMullvad(hardcoded, persisted models.MullvadServers) models.MullvadServers {
	if persisted.Timestamp <= hardcoded.Timestamp {
		return hardcoded
	}

	versionDiff := hardcoded.Version - persisted.Version
	if versionDiff > 0 {
		s.logVersionDiff("Mullvad", versionDiff)
		return hardcoded
	}

	s.logger.Info("Using Mullvad servers from file (%s more recent)",
		getUnixTimeDifference(persisted.Timestamp, hardcoded.Timestamp))
	return persisted
}

func (s *storage) mergeNordVPN(hardcoded, persisted models.NordvpnServers) models.NordvpnServers {
	if persisted.Timestamp <= hardcoded.Timestamp {
		return hardcoded
	}

	versionDiff := hardcoded.Version - persisted.Version
	if versionDiff > 0 {
		s.logVersionDiff("NordVPN", versionDiff)
		return hardcoded
	}

	s.logger.Info("Using NordVPN servers from file (%s more recent)",
		getUnixTimeDifference(persisted.Timestamp, hardcoded.Timestamp))
	return persisted
}

func (s *storage) mergePrivado(hardcoded, persisted models.PrivadoServers) models.PrivadoServers {
	if persisted.Timestamp <= hardcoded.Timestamp {
		return hardcoded
	}
	versionDiff := hardcoded.Version - persisted.Version
	if versionDiff > 0 {
		s.logVersionDiff("Privado", versionDiff)
		return hardcoded
	}
	s.logger.Info("Using Privado servers from file (%s more recent)",
		getUnixTimeDifference(persisted.Timestamp, hardcoded.Timestamp))
	return persisted
}

func (s *storage) mergePIA(hardcoded, persisted models.PiaServers) models.PiaServers {
	if persisted.Timestamp <= hardcoded.Timestamp {
		return hardcoded
	}
	versionDiff := hardcoded.Version - persisted.Version
	if versionDiff > 0 {
		s.logVersionDiff("Private Internet Access", versionDiff)
		return hardcoded
	}
	s.logger.Info("Using PIA servers from file (%s more recent)",
		getUnixTimeDifference(persisted.Timestamp, hardcoded.Timestamp))
	return persisted
}

func (s *storage) mergePrivatevpn(hardcoded, persisted models.PrivatevpnServers) models.PrivatevpnServers {
	if persisted.Timestamp <= hardcoded.Timestamp {
		return hardcoded
	}
	versionDiff := hardcoded.Version - persisted.Version
	if versionDiff > 0 {
		s.logVersionDiff("PrivateVPN", versionDiff)
		return hardcoded
	}
	s.logger.Info("Using Privatevpn servers from file (%s more recent)",
		getUnixTimeDifference(persisted.Timestamp, hardcoded.Timestamp))
	return persisted
}

func (s *storage) mergeProtonvpn(hardcoded, persisted models.ProtonvpnServers) models.ProtonvpnServers {
	if persisted.Timestamp <= hardcoded.Timestamp {
		return hardcoded
	}
	versionDiff := hardcoded.Version - persisted.Version
	if versionDiff > 0 {
		s.logVersionDiff("ProtonVPN", versionDiff)
		return hardcoded
	}
	s.logger.Info("Using Protonvpn servers from file (%s more recent)",
		getUnixTimeDifference(persisted.Timestamp, hardcoded.Timestamp))
	return persisted
}

func (s *storage) mergePureVPN(hardcoded, persisted models.PurevpnServers) models.PurevpnServers {
	if persisted.Timestamp <= hardcoded.Timestamp {
		return hardcoded
	}

	versionDiff := hardcoded.Version - persisted.Version
	if versionDiff > 0 {
		s.logVersionDiff("PureVPN", versionDiff)
		return hardcoded
	}

	s.logger.Info("Using PureVPN servers from file (%s more recent)",
		getUnixTimeDifference(persisted.Timestamp, hardcoded.Timestamp))
	return persisted
}

func (s *storage) mergeSurfshark(hardcoded, persisted models.SurfsharkServers) models.SurfsharkServers {
	if persisted.Timestamp <= hardcoded.Timestamp {
		return hardcoded
	}

	versionDiff := hardcoded.Version - persisted.Version
	if versionDiff > 0 {
		s.logVersionDiff("Surfshark", versionDiff)
		return hardcoded
	}

	s.logger.Info("Using Surfshark servers from file (%s more recent)",
		getUnixTimeDifference(persisted.Timestamp, hardcoded.Timestamp))
	return persisted
}

func (s *storage) mergeTorguard(hardcoded, persisted models.TorguardServers) models.TorguardServers {
	if persisted.Timestamp <= hardcoded.Timestamp {
		return hardcoded
	}
	versionDiff := hardcoded.Version - persisted.Version
	if versionDiff > 0 {
		s.logVersionDiff("Torguard", versionDiff)
		return hardcoded
	}
	s.logger.Info("Using Torguard servers from file (%s more recent)",
		getUnixTimeDifference(persisted.Timestamp, hardcoded.Timestamp))
	return persisted
}

func (s *storage) mergeVyprvpn(hardcoded, persisted models.VyprvpnServers) models.VyprvpnServers {
	if persisted.Timestamp <= hardcoded.Timestamp {
		return hardcoded
	}

	versionDiff := hardcoded.Version - persisted.Version
	if versionDiff > 0 {
		s.logVersionDiff("VyprVPN", versionDiff)
		return hardcoded
	}

	s.logger.Info("Using VyprVPN servers from file (%s more recent)",
		getUnixTimeDifference(persisted.Timestamp, hardcoded.Timestamp))
	return persisted
}

func (s *storage) mergeWindscribe(hardcoded, persisted models.WindscribeServers) models.WindscribeServers {
	if persisted.Timestamp <= hardcoded.Timestamp {
		return hardcoded
	}
	versionDiff := hardcoded.Version - persisted.Version
	if versionDiff > 0 {
		s.logVersionDiff("Windscribe", versionDiff)
		return hardcoded
	}
	s.logger.Info("Using Windscribe servers from file (%s more recent)",
		getUnixTimeDifference(persisted.Timestamp, hardcoded.Timestamp))
	return persisted
}
