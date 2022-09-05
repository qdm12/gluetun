package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gluetun/internal/models"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// readFromFile reads the servers from server.json.
// It only reads servers that have the same version as the hardcoded servers version
// to avoid JSON unmarshaling errors.
func (s *Storage) readFromFile(filepath string, hardcodedVersions map[string]uint16) (
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

	return s.extractServersFromBytes(b, hardcodedVersions)
}

func (s *Storage) extractServersFromBytes(b []byte, hardcodedVersions map[string]uint16) (
	servers models.AllServers, err error) {
	rawMessages := make(map[string]json.RawMessage)
	if err := json.Unmarshal(b, &rawMessages); err != nil {
		return servers, fmt.Errorf("cannot decode servers: %w", err)
	}

	// Note schema version is at map key "version" as number

	allProviders := providers.All()
	servers.ProviderToServers = make(map[string]models.Servers, len(allProviders))
	titleCaser := cases.Title(language.English)
	for _, provider := range allProviders {
		hardcodedVersion, ok := hardcodedVersions[provider]
		if !ok {
			panic(fmt.Sprintf("provider %s not found in hardcoded servers map; "+
				"did you add the provider key in the embedded servers.json?", provider))
		}

		rawMessage, ok := rawMessages[provider]
		if !ok {
			// If the provider is not found in the data bytes, just don't set it in
			// the providers map. That way the hardcoded servers will override them.
			// This is user provided and could come from different sources in the
			// future (e.g. a file or API request).
			continue
		}

		mergedServers, versionsMatch, err := s.readServers(provider,
			hardcodedVersion, rawMessage, titleCaser)
		if err != nil {
			return models.AllServers{}, err
		} else if !versionsMatch {
			// mergedServers is the empty struct in this case, so don't set the key
			// in the providerToServers map.
			continue
		}
		servers.ProviderToServers[provider] = mergedServers
	}

	return servers, nil
}

func (s *Storage) readServers(provider string, hardcodedVersion uint16,
	rawMessage json.RawMessage, titleCaser cases.Caser) (servers models.Servers,
	versionsMatch bool, err error) {
	provider = titleCaser.String(provider)

	var versionObject struct {
		Version uint16 `json:"version"`
	}

	err = json.Unmarshal(rawMessage, &versionObject)
	if err != nil {
		return servers, false, fmt.Errorf("cannot decode servers version for provider %s: %w",
			provider, err)
	}

	persistedVersion := versionObject.Version

	versionsMatch = hardcodedVersion == persistedVersion
	if !versionsMatch {
		s.logger.Info(fmt.Sprintf(
			"%s servers from file discarded because they have "+
				"version %d and hardcoded servers have version %d",
			provider, persistedVersion, hardcodedVersion))
		return servers, versionsMatch, nil
	}

	err = json.Unmarshal(rawMessage, &servers)
	if err != nil {
		return servers, false, fmt.Errorf("cannot decode servers for provider %s: %w",
			provider, err)
	}

	return servers, versionsMatch, nil
}
