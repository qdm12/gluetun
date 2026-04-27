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

// readFromFile reads the servers data starting from the given manifest file path.
// It only reads servers that have the same version as the hardcoded servers version
// to avoid JSON decoding errors.
func (s *Storage) readFromFile(manifestPath string, hardcodedVersions map[string]uint16) (
	servers models.AllServers, found bool, err error,
) {
	file, err := os.Open(manifestPath)
	if os.IsNotExist(err) {
		return servers, false, nil
	} else if err != nil {
		return servers, false, err
	}

	b, err := io.ReadAll(file)
	if err != nil {
		return servers, true, err
	}

	if err := file.Close(); err != nil {
		return servers, true, err
	}

	servers, err = s.extractServersFromBytes(b, hardcodedVersions)
	return servers, true, err
}

func (s *Storage) extractServersFromBytes(b []byte, hardcodedVersions map[string]uint16) (
	servers models.AllServers, err error,
) {
	rawMessages := make(map[string]json.RawMessage)
	if err := json.Unmarshal(b, &rawMessages); err != nil {
		return servers, fmt.Errorf("decoding servers: %w", err)
	}

	// Note schema version is at map key "version" as number
	if rawVersion, ok := rawMessages["version"]; ok {
		err := json.Unmarshal(rawVersion, &servers.Version)
		if err != nil {
			return servers, fmt.Errorf("decoding servers schema version: %w", err)
		}
	}

	allProviders := providers.All()
	servers.ProviderToServers = make(map[string]models.Servers, len(allProviders))
	titleCaser := cases.Title(language.English)
	for _, provider := range allProviders {
		hardcodedVersion, ok := hardcodedVersions[provider]
		if !ok {
			panicOnProviderMissingHardcoded(provider)
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
	versionsMatch bool, err error,
) {
	provider = titleCaser.String(provider)

	var metadata struct {
		Version   uint16 `json:"version"`
		Timestamp int64  `json:"timestamp"`
		Filepath  string `json:"filepath"`
	}

	err = json.Unmarshal(rawMessage, &metadata)
	if err != nil {
		return servers, false, fmt.Errorf("decoding servers version for provider %s: %w",
			provider, err)
	}

	if metadata.Filepath != "" {
		providerFile, err := os.Open(metadata.Filepath)
		if os.IsNotExist(err) {
			return models.Servers{}, false, nil
		} else if err != nil {
			return models.Servers{}, false, fmt.Errorf("opening servers file %s for provider %s: %w",
				metadata.Filepath, provider, err)
		}
		defer providerFile.Close()

		var referencedServers models.Servers
		err = json.NewDecoder(providerFile).Decode(&referencedServers)
		if err != nil {
			return models.Servers{}, false, fmt.Errorf("decoding servers file %s for provider %s: %w",
				metadata.Filepath, provider, err)
		}

		versionsMatch = referencedServers.Version == hardcodedVersion
		if !versionsMatch {
			if referencedServers.Preferred {
				s.logger.Warn(fmt.Sprintf(
					"%s preferred servers from file %s discarded because they have version %d and hardcoded servers have version %d",
					provider, metadata.Filepath, referencedServers.Version, hardcodedVersion))
			} else {
				s.logger.Info(fmt.Sprintf(
					"%s servers from file %s discarded because they have version %d and hardcoded servers have version %d",
					provider, metadata.Filepath, referencedServers.Version, hardcodedVersion))
			}
			return models.Servers{}, false, nil
		}

		referencedServers.Filepath = metadata.Filepath
		return referencedServers, true, nil
	}

	err = json.Unmarshal(rawMessage, &servers)
	if err != nil {
		return servers, false, fmt.Errorf("decoding servers for provider %s: %w",
			provider, err)
	}

	versionsMatch = servers.Version == hardcodedVersion
	if !versionsMatch {
		if servers.Preferred {
			s.logger.Warn(fmt.Sprintf(
				"%s preferred servers from file discarded because they have version %d and hardcoded servers have version %d",
				provider, servers.Version, hardcodedVersion))
		} else {
			s.logger.Info(fmt.Sprintf(
				"%s servers from file discarded because they have version %d and hardcoded servers have version %d",
				provider, servers.Version, hardcodedVersion))
		}
		return servers, false, nil
	}

	return servers, true, nil
}
