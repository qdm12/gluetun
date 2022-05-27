package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/qdm12/gluetun/internal/constants/providers"
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

func (s *Storage) extractServersFromBytes(b []byte, hardcoded models.AllServers) (
	servers models.AllServers, err error) {
	var versions allVersions
	if err := json.Unmarshal(b, &versions); err != nil {
		return servers, fmt.Errorf("cannot decode versions: %w", err)
	}

	var rawMessages allJSONRawMessages
	if err := json.Unmarshal(b, &rawMessages); err != nil {
		return servers, fmt.Errorf("cannot decode servers: %w", err)
	}

	type element struct {
		provider      string
		hardcoded     models.Servers
		serverVersion serverVersion
		rawMessage    json.RawMessage
		target        *models.Servers
	}
	elements := []element{
		{
			provider:      providers.Cyberghost,
			hardcoded:     hardcoded.Cyberghost,
			serverVersion: versions.Cyberghost,
			rawMessage:    rawMessages.Cyberghost,
			target:        &servers.Cyberghost,
		},
		{
			provider:      providers.Expressvpn,
			hardcoded:     hardcoded.Expressvpn,
			serverVersion: versions.Expressvpn,
			rawMessage:    rawMessages.Expressvpn,
			target:        &servers.Expressvpn,
		},
		{
			provider:      providers.Fastestvpn,
			hardcoded:     hardcoded.Fastestvpn,
			serverVersion: versions.Fastestvpn,
			rawMessage:    rawMessages.Fastestvpn,
			target:        &servers.Fastestvpn,
		},
		{
			provider:      providers.HideMyAss,
			hardcoded:     hardcoded.HideMyAss,
			serverVersion: versions.HideMyAss,
			rawMessage:    rawMessages.HideMyAss,
			target:        &servers.HideMyAss,
		},
		{
			provider:      providers.Ipvanish,
			hardcoded:     hardcoded.Ipvanish,
			serverVersion: versions.Ipvanish,
			rawMessage:    rawMessages.Ipvanish,
			target:        &servers.Ipvanish,
		},
		{
			provider:      providers.Ivpn,
			hardcoded:     hardcoded.Ivpn,
			serverVersion: versions.Ivpn,
			rawMessage:    rawMessages.Ivpn,
			target:        &servers.Ivpn,
		},
		{
			provider:      providers.Mullvad,
			hardcoded:     hardcoded.Mullvad,
			serverVersion: versions.Mullvad,
			rawMessage:    rawMessages.Mullvad,
			target:        &servers.Mullvad,
		},
		{
			provider:      providers.Nordvpn,
			hardcoded:     hardcoded.Nordvpn,
			serverVersion: versions.Nordvpn,
			rawMessage:    rawMessages.Nordvpn,
			target:        &servers.Nordvpn,
		},
		{
			provider:      providers.Perfectprivacy,
			hardcoded:     hardcoded.Perfectprivacy,
			serverVersion: versions.Perfectprivacy,
			rawMessage:    rawMessages.Perfectprivacy,
			target:        &servers.Perfectprivacy,
		},
		{
			provider:      providers.Privado,
			hardcoded:     hardcoded.Privado,
			serverVersion: versions.Privado,
			rawMessage:    rawMessages.Privado,
			target:        &servers.Privado,
		},
		{
			provider:      providers.PrivateInternetAccess,
			hardcoded:     hardcoded.Pia,
			serverVersion: versions.Pia,
			rawMessage:    rawMessages.Pia,
			target:        &servers.Pia,
		},
		{
			provider:      providers.Privatevpn,
			hardcoded:     hardcoded.Privatevpn,
			serverVersion: versions.Privatevpn,
			rawMessage:    rawMessages.Privatevpn,
			target:        &servers.Privatevpn,
		},
		{
			provider:      providers.Protonvpn,
			hardcoded:     hardcoded.Protonvpn,
			serverVersion: versions.Protonvpn,
			rawMessage:    rawMessages.Protonvpn,
			target:        &servers.Protonvpn,
		},
		{
			provider:      providers.Purevpn,
			hardcoded:     hardcoded.Purevpn,
			serverVersion: versions.Purevpn,
			rawMessage:    rawMessages.Purevpn,
			target:        &servers.Purevpn,
		},
		{
			provider:      providers.Surfshark,
			hardcoded:     hardcoded.Surfshark,
			serverVersion: versions.Surfshark,
			rawMessage:    rawMessages.Surfshark,
			target:        &servers.Surfshark,
		},
		{
			provider:      providers.Torguard,
			hardcoded:     hardcoded.Torguard,
			serverVersion: versions.Torguard,
			rawMessage:    rawMessages.Torguard,
			target:        &servers.Torguard,
		},
		{
			provider:      providers.VPNUnlimited,
			hardcoded:     hardcoded.VPNUnlimited,
			serverVersion: versions.VPNUnlimited,
			rawMessage:    rawMessages.VPNUnlimited,
			target:        &servers.VPNUnlimited,
		},
		{
			provider:      providers.Vyprvpn,
			hardcoded:     hardcoded.Vyprvpn,
			serverVersion: versions.Vyprvpn,
			rawMessage:    rawMessages.Vyprvpn,
			target:        &servers.Vyprvpn,
		},
		{
			provider:      providers.Wevpn,
			hardcoded:     hardcoded.Wevpn,
			serverVersion: versions.Wevpn,
			rawMessage:    rawMessages.Wevpn,
			target:        &servers.Wevpn,
		},
		{
			provider:      providers.Windscribe,
			hardcoded:     hardcoded.Windscribe,
			serverVersion: versions.Windscribe,
			rawMessage:    rawMessages.Windscribe,
			target:        &servers.Windscribe,
		},
	}

	for _, element := range elements {
		*element.target, err = s.readServers(element.provider,
			element.hardcoded, element.serverVersion, element.rawMessage)
		if err != nil {
			return servers, err
		}
	}

	return servers, nil
}

var (
	errDecodeProvider = errors.New("cannot decode servers for provider")
)

func (s *Storage) readServers(provider string, hardcoded models.Servers,
	serverVersion serverVersion, rawMessage json.RawMessage) (
	servers models.Servers, err error) {
	provider = strings.Title(provider)
	if hardcoded.Version != serverVersion.Version {
		s.logVersionDiff(provider, hardcoded.Version, serverVersion.Version)
		return servers, nil
	}

	err = json.Unmarshal(rawMessage, &servers)
	if err != nil {
		return servers, fmt.Errorf("%w: %s: %s", errDecodeProvider, provider, err)
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
	Perfectprivacy serverVersion `json:"perfect privacy"`
	Privado        serverVersion `json:"privado"`
	Pia            serverVersion `json:"private internet access"`
	Privatevpn     serverVersion `json:"privatevpn"`
	Protonvpn      serverVersion `json:"protonvpn"`
	Purevpn        serverVersion `json:"purevpn"`
	Surfshark      serverVersion `json:"surfshark"`
	Torguard       serverVersion `json:"torguard"`
	VPNUnlimited   serverVersion `json:"vpn unlimited"`
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
	Perfectprivacy json.RawMessage `json:"perfect privacy"`
	Privado        json.RawMessage `json:"privado"`
	Pia            json.RawMessage `json:"private internet access"`
	Privatevpn     json.RawMessage `json:"privatevpn"`
	Protonvpn      json.RawMessage `json:"protonvpn"`
	Purevpn        json.RawMessage `json:"purevpn"`
	Surfshark      json.RawMessage `json:"surfshark"`
	Torguard       json.RawMessage `json:"torguard"`
	VPNUnlimited   json.RawMessage `json:"vpn unlimited"`
	Vyprvpn        json.RawMessage `json:"vyprvpn"`
	Wevpn          json.RawMessage `json:"wevpn"`
	Windscribe     json.RawMessage `json:"windscribe"`
}
