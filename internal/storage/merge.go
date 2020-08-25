package storage

import "github.com/qdm12/gluetun/internal/models"

func mergeServers(hardcoded, persistent models.AllServers) (merged models.AllServers) {
	merged.Version = hardcoded.Version
	merged.Cyberghost = hardcoded.Cyberghost
	if persistent.Cyberghost.Timestamp > hardcoded.Cyberghost.Timestamp {
		merged.Cyberghost = persistent.Cyberghost
	}
	merged.Mullvad = hardcoded.Mullvad
	if persistent.Mullvad.Timestamp > hardcoded.Mullvad.Timestamp {
		merged.Mullvad = persistent.Mullvad
	}
	merged.Nordvpn = hardcoded.Nordvpn
	if persistent.Nordvpn.Timestamp > hardcoded.Nordvpn.Timestamp {
		merged.Nordvpn = persistent.Nordvpn
	}
	merged.Pia = hardcoded.Pia
	if persistent.Pia.Timestamp > hardcoded.Pia.Timestamp {
		merged.Pia = persistent.Pia
	}
	merged.Purevpn = hardcoded.Purevpn
	if persistent.Purevpn.Timestamp > hardcoded.Purevpn.Timestamp {
		merged.Purevpn = persistent.Purevpn
	}
	merged.Surfshark = hardcoded.Surfshark
	if persistent.Surfshark.Timestamp > hardcoded.Surfshark.Timestamp {
		merged.Surfshark = persistent.Surfshark
	}
	merged.Vyprvpn = hardcoded.Vyprvpn
	if persistent.Vyprvpn.Timestamp > hardcoded.Vyprvpn.Timestamp {
		merged.Vyprvpn = persistent.Vyprvpn
	}
	merged.Windscribe = hardcoded.Windscribe
	if persistent.Windscribe.Timestamp > hardcoded.Windscribe.Timestamp {
		merged.Windscribe = persistent.Windscribe
	}
	return merged
}
