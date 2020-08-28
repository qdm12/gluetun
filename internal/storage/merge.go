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

func (s *storage) mergeServers(hardcoded, persistent models.AllServers) (merged models.AllServers) {
	merged.Version = hardcoded.Version
	merged.Cyberghost = hardcoded.Cyberghost
	if persistent.Cyberghost.Timestamp > hardcoded.Cyberghost.Timestamp {
		s.logger.Info("Using Cyberghost servers from file (%s more recent)",
			getUnixTimeDifference(persistent.Cyberghost.Timestamp, hardcoded.Cyberghost.Timestamp))
		merged.Cyberghost = persistent.Cyberghost
	}
	merged.Mullvad = hardcoded.Mullvad
	if persistent.Mullvad.Timestamp > hardcoded.Mullvad.Timestamp {
		s.logger.Info("Using Mullvad servers from file (%s more recent)",
			getUnixTimeDifference(persistent.Mullvad.Timestamp, hardcoded.Mullvad.Timestamp))
		merged.Mullvad = persistent.Mullvad
	}
	merged.Nordvpn = hardcoded.Nordvpn
	if persistent.Nordvpn.Timestamp > hardcoded.Nordvpn.Timestamp {
		s.logger.Info("Using Nordvpn servers from file (%s more recent)",
			getUnixTimeDifference(persistent.Nordvpn.Timestamp, hardcoded.Nordvpn.Timestamp))
		merged.Nordvpn = persistent.Nordvpn
	}
	merged.Pia = hardcoded.Pia
	if persistent.Pia.Timestamp > hardcoded.Pia.Timestamp {
		s.logger.Info("Using Private Internet Access servers from file (%s more recent)",
			getUnixTimeDifference(persistent.Pia.Timestamp, hardcoded.Pia.Timestamp))
		merged.Pia = persistent.Pia
	}
	merged.PiaOld = hardcoded.PiaOld
	if persistent.PiaOld.Timestamp > hardcoded.PiaOld.Timestamp {
		s.logger.Info("Using Private Internet Access older servers from file (%s more recent)",
			getUnixTimeDifference(persistent.PiaOld.Timestamp, hardcoded.PiaOld.Timestamp))
		merged.PiaOld = persistent.PiaOld
	}
	merged.Purevpn = hardcoded.Purevpn
	if persistent.Purevpn.Timestamp > hardcoded.Purevpn.Timestamp {
		s.logger.Info("Using Purevpn servers from file (%s more recent)",
			getUnixTimeDifference(persistent.Purevpn.Timestamp, hardcoded.Purevpn.Timestamp))
		merged.Purevpn = persistent.Purevpn
	}
	merged.Surfshark = hardcoded.Surfshark
	if persistent.Surfshark.Timestamp > hardcoded.Surfshark.Timestamp {
		s.logger.Info("Using Surfshark servers from file (%s more recent)",
			getUnixTimeDifference(persistent.Surfshark.Timestamp, hardcoded.Surfshark.Timestamp))
		merged.Surfshark = persistent.Surfshark
	}
	merged.Vyprvpn = hardcoded.Vyprvpn
	if persistent.Vyprvpn.Timestamp > hardcoded.Vyprvpn.Timestamp {
		s.logger.Info("Using Vyprvpn servers from file (%s more recent)",
			getUnixTimeDifference(persistent.Vyprvpn.Timestamp, hardcoded.Vyprvpn.Timestamp))
		merged.Vyprvpn = persistent.Vyprvpn
	}
	merged.Windscribe = hardcoded.Windscribe
	if persistent.Windscribe.Timestamp > hardcoded.Windscribe.Timestamp {
		s.logger.Info("Using Windscribe servers from file (%s more recent)",
			getUnixTimeDifference(persistent.Windscribe.Timestamp, hardcoded.Windscribe.Timestamp))
		merged.Windscribe = persistent.Windscribe
	}
	return merged
}
