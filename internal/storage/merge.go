package storage

import (
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

func (s *storage) mergeServers(hardcoded, persisted models.AllServers) models.AllServers {
	return models.AllServers{
		Version:    hardcoded.Version,
		Cyberghost: s.mergeCyberghost(hardcoded.Cyberghost, persisted.Cyberghost),
		Mullvad:    s.mergeMullvad(hardcoded.Mullvad, persisted.Mullvad),
		Nordvpn:    s.mergeNordVPN(hardcoded.Nordvpn, persisted.Nordvpn),
		Pia:        s.mergePIA(hardcoded.Pia, persisted.Pia),
		Privado:    s.mergePrivado(hardcoded.Privado, persisted.Privado),
		Purevpn:    s.mergePureVPN(hardcoded.Purevpn, persisted.Purevpn),
		Surfshark:  s.mergeSurfshark(hardcoded.Surfshark, persisted.Surfshark),
		Vyprvpn:    s.mergeVyprvpn(hardcoded.Vyprvpn, persisted.Vyprvpn),
		Windscribe: s.mergeWindscribe(hardcoded.Windscribe, persisted.Windscribe),
	}
}

func (s *storage) mergeCyberghost(hardcoded, persisted models.CyberghostServers) models.CyberghostServers {
	if persisted.Timestamp <= hardcoded.Timestamp {
		return hardcoded
	}
	s.logger.Info("Using Cyberghost servers from file (%s more recent)",
		getUnixTimeDifference(persisted.Timestamp, hardcoded.Timestamp))
	return persisted
}

func (s *storage) mergeMullvad(hardcoded, persisted models.MullvadServers) models.MullvadServers {
	if persisted.Timestamp <= hardcoded.Timestamp {
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
	s.logger.Info("Using NordVPN servers from file (%s more recent)",
		getUnixTimeDifference(persisted.Timestamp, hardcoded.Timestamp))
	return persisted
}

func (s *storage) mergePIA(hardcoded, persisted models.PiaServers) models.PiaServers {
	if persisted.Timestamp <= hardcoded.Timestamp {
		return hardcoded
	}
	versionDiff := hardcoded.Version - persisted.Version
	if versionDiff > 0 {
		s.logger.Info(
			"PIA servers from file discarded because they are %d versions behind",
			versionDiff)
		return hardcoded
	}
	s.logger.Info("Using PIA servers from file (%s more recent)",
		getUnixTimeDifference(persisted.Timestamp, hardcoded.Timestamp))
	return persisted
}

func (s *storage) mergePrivado(hardcoded, persisted models.PrivadoServers) models.PrivadoServers {
	if persisted.Timestamp <= hardcoded.Timestamp {
		return hardcoded
	}
	versionDiff := hardcoded.Version - persisted.Version
	if versionDiff > 0 {
		s.logger.Info(
			"Privado servers from file discarded because they are %d versions behind",
			versionDiff)
		return hardcoded
	}
	s.logger.Info("Using Privado servers from file (%s more recent)",
		getUnixTimeDifference(persisted.Timestamp, hardcoded.Timestamp))
	return persisted
}

func (s *storage) mergePureVPN(hardcoded, persisted models.PurevpnServers) models.PurevpnServers {
	if persisted.Timestamp <= hardcoded.Timestamp {
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
	s.logger.Info("Using Surfshark servers from file (%s more recent)",
		getUnixTimeDifference(persisted.Timestamp, hardcoded.Timestamp))
	return persisted
}

func (s *storage) mergeVyprvpn(hardcoded, persisted models.VyprvpnServers) models.VyprvpnServers {
	if persisted.Timestamp <= hardcoded.Timestamp {
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
		s.logger.Info(
			"Windscribe servers from file discarded because they are %d versions behind",
			versionDiff)
		return hardcoded
	}
	s.logger.Info("Using Windscribe servers from file (%s more recent)",
		getUnixTimeDifference(persisted.Timestamp, hardcoded.Timestamp))
	return persisted
}
