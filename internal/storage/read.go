package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/qdm12/gluetun/internal/models"
)

// readFromFile reads the servers from server.json.
// It only reads servers that have the same version as the hardcoded servers version
// to avoid JSON unmarshaling errors.
func (s *Storage) readFromFile(filepath string, hardcoded models.AllServers) (
	servers models.AllServers, err error) {
	file, err := os.Open(filepath)
	if os.IsNotExist(err) {
		return servers, nil
	} else if err != nil {
		return servers, err
	}

	b, err := io.ReadAll(file)
	if err != nil {
		return servers, err
	}

	if err := file.Close(); err != nil {
		return servers, err
	}

	return s.extractServersFromBytes(b, hardcoded)
}

var (
	errDecodeProvider = errors.New("cannot decode servers for provider")
)

func (s *Storage) extractServersFromBytes(b []byte, hardcoded models.AllServers) ( //nolint:gocognit,gocyclo
	servers models.AllServers, err error) {
	var versions allVersions
	if err := json.Unmarshal(b, &versions); err != nil {
		return servers, fmt.Errorf("cannot decode versions: %w", err)
	}

	var rawMessages allJSONRawMessages
	if err := json.Unmarshal(b, &rawMessages); err != nil {
		return servers, fmt.Errorf("cannot decode servers: %w", err)
	}

	// TODO simplify with generics in Go 1.18

	if hardcoded.Cyberghost.Version != versions.Cyberghost.Version {
		s.logVersionDiff("Cyberghost", hardcoded.Cyberghost.Version, versions.Cyberghost.Version)
	} else {
		err = json.Unmarshal(rawMessages.Cyberghost, &servers.Cyberghost)
		if err != nil {
			return servers, fmt.Errorf("%w: %s: %s", errDecodeProvider, "Cyberghost", err)
		}
	}

	if hardcoded.Expressvpn.Version != versions.Expressvpn.Version {
		s.logVersionDiff("Expressvpn", hardcoded.Expressvpn.Version, versions.Expressvpn.Version)
	} else {
		err = json.Unmarshal(rawMessages.Expressvpn, &servers.Expressvpn)
		if err != nil {
			return servers, fmt.Errorf("%w: %s: %s", errDecodeProvider, "Expressvpn", err)
		}
	}

	if hardcoded.Fastestvpn.Version != versions.Fastestvpn.Version {
		s.logVersionDiff("Fastestvpn", hardcoded.Fastestvpn.Version, versions.Fastestvpn.Version)
	} else {
		err = json.Unmarshal(rawMessages.Fastestvpn, &servers.Fastestvpn)
		if err != nil {
			return servers, fmt.Errorf("%w: %s: %s", errDecodeProvider, "Fastestvpn", err)
		}
	}

	if hardcoded.HideMyAss.Version != versions.HideMyAss.Version {
		s.logVersionDiff("HideMyAss", hardcoded.HideMyAss.Version, versions.HideMyAss.Version)
	} else {
		err = json.Unmarshal(rawMessages.HideMyAss, &servers.HideMyAss)
		if err != nil {
			return servers, fmt.Errorf("%w: %s: %s", errDecodeProvider, "HideMyAss", err)
		}
	}

	if hardcoded.Ipvanish.Version != versions.Ipvanish.Version {
		s.logVersionDiff("Ipvanish", hardcoded.Ipvanish.Version, versions.Ipvanish.Version)
	} else {
		err = json.Unmarshal(rawMessages.Ipvanish, &servers.Ipvanish)
		if err != nil {
			return servers, fmt.Errorf("%w: %s: %s", errDecodeProvider, "Ipvanish", err)
		}
	}

	if hardcoded.Ivpn.Version != versions.Ivpn.Version {
		s.logVersionDiff("Ivpn", hardcoded.Ivpn.Version, versions.Ivpn.Version)
	} else {
		err = json.Unmarshal(rawMessages.Ivpn, &servers.Ivpn)
		if err != nil {
			return servers, fmt.Errorf("%w: %s: %s", errDecodeProvider, "Ivpn", err)
		}
	}

	if hardcoded.Mullvad.Version != versions.Mullvad.Version {
		s.logVersionDiff("Mullvad", hardcoded.Mullvad.Version, versions.Mullvad.Version)
	} else {
		err = json.Unmarshal(rawMessages.Mullvad, &servers.Mullvad)
		if err != nil {
			return servers, fmt.Errorf("%w: %s: %s", errDecodeProvider, "Mullvad", err)
		}
	}

	if hardcoded.Nordvpn.Version != versions.Nordvpn.Version {
		s.logVersionDiff("Nordvpn", hardcoded.Nordvpn.Version, versions.Nordvpn.Version)
	} else {
		err = json.Unmarshal(rawMessages.Nordvpn, &servers.Nordvpn)
		if err != nil {
			return servers, fmt.Errorf("%w: %s: %s", errDecodeProvider, "Nordvpn", err)
		}
	}

	if hardcoded.Perfectprivacy.Version != versions.Perfectprivacy.Version {
		s.logVersionDiff("Perfect Privacy", hardcoded.Perfectprivacy.Version, versions.Perfectprivacy.Version)
	} else {
		err = json.Unmarshal(rawMessages.Perfectprivacy, &servers.Perfectprivacy)
		if err != nil {
			return servers, fmt.Errorf("%w: %s: %s", errDecodeProvider, "Perfect Privacy", err)
		}
	}

	if hardcoded.Privado.Version != versions.Privado.Version {
		s.logVersionDiff("Privado", hardcoded.Privado.Version, versions.Privado.Version)
	} else {
		err = json.Unmarshal(rawMessages.Privado, &servers.Privado)
		if err != nil {
			return servers, fmt.Errorf("%w: %s: %s", errDecodeProvider, "Privado", err)
		}
	}

	if hardcoded.Pia.Version != versions.Pia.Version {
		s.logVersionDiff("Pia", hardcoded.Pia.Version, versions.Pia.Version)
	} else {
		err = json.Unmarshal(rawMessages.Pia, &servers.Pia)
		if err != nil {
			return servers, fmt.Errorf("%w: %s: %s", errDecodeProvider, "Pia", err)
		}
	}

	if hardcoded.Privatevpn.Version != versions.Privatevpn.Version {
		s.logVersionDiff("Privatevpn", hardcoded.Privatevpn.Version, versions.Privatevpn.Version)
	} else {
		err = json.Unmarshal(rawMessages.Privatevpn, &servers.Privatevpn)
		if err != nil {
			return servers, fmt.Errorf("%w: %s: %s", errDecodeProvider, "Privatevpn", err)
		}
	}

	if hardcoded.Protonvpn.Version != versions.Protonvpn.Version {
		s.logVersionDiff("Protonvpn", hardcoded.Protonvpn.Version, versions.Protonvpn.Version)
	} else {
		err = json.Unmarshal(rawMessages.Protonvpn, &servers.Protonvpn)
		if err != nil {
			return servers, fmt.Errorf("%w: %s: %s", errDecodeProvider, "Protonvpn", err)
		}
	}

	if hardcoded.Purevpn.Version != versions.Purevpn.Version {
		s.logVersionDiff("Purevpn", hardcoded.Purevpn.Version, versions.Purevpn.Version)
	} else {
		err = json.Unmarshal(rawMessages.Purevpn, &servers.Purevpn)
		if err != nil {
			return servers, fmt.Errorf("%w: %s: %s", errDecodeProvider, "Purevpn", err)
		}
	}

	if hardcoded.Surfshark.Version != versions.Surfshark.Version {
		s.logVersionDiff("Surfshark", hardcoded.Surfshark.Version, versions.Surfshark.Version)
	} else {
		err = json.Unmarshal(rawMessages.Surfshark, &servers.Surfshark)
		if err != nil {
			return servers, fmt.Errorf("%w: %s: %s", errDecodeProvider, "Surfshark", err)
		}
	}

	if hardcoded.Torguard.Version != versions.Torguard.Version {
		s.logVersionDiff("Torguard", hardcoded.Torguard.Version, versions.Torguard.Version)
	} else {
		err = json.Unmarshal(rawMessages.Torguard, &servers.Torguard)
		if err != nil {
			return servers, fmt.Errorf("%w: %s: %s", errDecodeProvider, "Torguard", err)
		}
	}

	if hardcoded.VPNUnlimited.Version != versions.VPNUnlimited.Version {
		s.logVersionDiff("VPNUnlimited", hardcoded.VPNUnlimited.Version, versions.VPNUnlimited.Version)
	} else {
		err = json.Unmarshal(rawMessages.VPNUnlimited, &servers.VPNUnlimited)
		if err != nil {
			return servers, fmt.Errorf("%w: %s: %s", errDecodeProvider, "VPNUnlimited", err)
		}
	}

	if hardcoded.Vyprvpn.Version != versions.Vyprvpn.Version {
		s.logVersionDiff("Vyprvpn", hardcoded.Vyprvpn.Version, versions.Vyprvpn.Version)
	} else {
		err = json.Unmarshal(rawMessages.Vyprvpn, &servers.Vyprvpn)
		if err != nil {
			return servers, fmt.Errorf("%w: %s: %s", errDecodeProvider, "Vyprvpn", err)
		}
	}

	if hardcoded.Wevpn.Version != versions.Wevpn.Version {
		s.logVersionDiff("Wevpn", hardcoded.Wevpn.Version, versions.Wevpn.Version)
	} else {
		err = json.Unmarshal(rawMessages.Wevpn, &servers.Wevpn)
		if err != nil {
			return servers, fmt.Errorf("%w: %s: %s", errDecodeProvider, "Wevpn", err)
		}
	}

	if hardcoded.Windscribe.Version != versions.Windscribe.Version {
		s.logVersionDiff("Windscribe", hardcoded.Windscribe.Version, versions.Windscribe.Version)
	} else {
		err = json.Unmarshal(rawMessages.Windscribe, &servers.Windscribe)
		if err != nil {
			return servers, fmt.Errorf("%w: %s: %s", errDecodeProvider, "Windscribe", err)
		}
	}

	return servers, nil
}

// allVersions is a subset of models.AllServers structure used to track
// versions to avoid unmarshaling errors.
type allVersions struct {
	Version        uint16        `json:"version"` // used for migration of the top level scheme
	Cyberghost     serverVersion `json:"cyberghost"`
	Expressvpn     serverVersion `json:"expressvpn"`
	Fastestvpn     serverVersion `json:"fastestvpn"`
	HideMyAss      serverVersion `json:"hidemyass"`
	Ipvanish       serverVersion `json:"ipvanish"`
	Ivpn           serverVersion `json:"ivpn"`
	Mullvad        serverVersion `json:"mullvad"`
	Nordvpn        serverVersion `json:"nordvpn"`
	Perfectprivacy serverVersion `json:"perfectprivacy"`
	Privado        serverVersion `json:"privado"`
	Pia            serverVersion `json:"pia"`
	Privatevpn     serverVersion `json:"privatevpn"`
	Protonvpn      serverVersion `json:"protonvpn"`
	Purevpn        serverVersion `json:"purevpn"`
	Surfshark      serverVersion `json:"surfshark"`
	Torguard       serverVersion `json:"torguard"`
	VPNUnlimited   serverVersion `json:"vpnunlimited"`
	Vyprvpn        serverVersion `json:"vyprvpn"`
	Wevpn          serverVersion `json:"wevpn"`
	Windscribe     serverVersion `json:"windscribe"`
}

type serverVersion struct {
	Version uint16 `json:"version"`
}

// allJSONRawMessages is to delay decoding of each provider servers.
type allJSONRawMessages struct {
	Version        uint16          `json:"version"` // used for migration of the top level scheme
	Cyberghost     json.RawMessage `json:"cyberghost"`
	Expressvpn     json.RawMessage `json:"expressvpn"`
	Fastestvpn     json.RawMessage `json:"fastestvpn"`
	HideMyAss      json.RawMessage `json:"hidemyass"`
	Ipvanish       json.RawMessage `json:"ipvanish"`
	Ivpn           json.RawMessage `json:"ivpn"`
	Mullvad        json.RawMessage `json:"mullvad"`
	Nordvpn        json.RawMessage `json:"nordvpn"`
	Perfectprivacy json.RawMessage `json:"perfectprivacy"`
	Privado        json.RawMessage `json:"privado"`
	Pia            json.RawMessage `json:"pia"`
	Privatevpn     json.RawMessage `json:"privatevpn"`
	Protonvpn      json.RawMessage `json:"protonvpn"`
	Purevpn        json.RawMessage `json:"purevpn"`
	Surfshark      json.RawMessage `json:"surfshark"`
	Torguard       json.RawMessage `json:"torguard"`
	VPNUnlimited   json.RawMessage `json:"vpnunlimited"`
	Vyprvpn        json.RawMessage `json:"vyprvpn"`
	Wevpn          json.RawMessage `json:"wevpn"`
	Windscribe     json.RawMessage `json:"windscribe"`
}
