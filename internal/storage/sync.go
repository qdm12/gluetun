package storage

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/os"
)

const (
	jsonFilepath = "/gluetun/servers.json"
)

func countServers(allServers models.AllServers) int {
	return len(allServers.Cyberghost.Servers) +
		len(allServers.Mullvad.Servers) +
		len(allServers.Nordvpn.Servers) +
		len(allServers.Pia.Servers) +
		len(allServers.Privado.Servers) +
		len(allServers.Purevpn.Servers) +
		len(allServers.Surfshark.Servers) +
		len(allServers.Vyprvpn.Servers) +
		len(allServers.Windscribe.Servers)
}

func (s *storage) SyncServers(hardcodedServers models.AllServers, write bool) (
	allServers models.AllServers, err error) {
	// Eventually read file
	var serversOnFile models.AllServers
	file, err := s.os.OpenFile(jsonFilepath, os.O_RDONLY, 0)
	if err != nil && !os.IsNotExist(err) {
		return allServers, err
	}
	if err == nil {
		var serversOnFile models.AllServers
		decoder := json.NewDecoder(file)
		if err := decoder.Decode(&serversOnFile); err != nil {
			_ = file.Close()
			return allServers, err
		}
		return allServers, file.Close()
	}

	// Merge data from file and hardcoded
	s.logger.Info("Merging by most recent %d hardcoded servers and %d servers read from %s",
		countServers(hardcodedServers), countServers(serversOnFile), jsonFilepath)
	allServers = s.mergeServers(hardcodedServers, serversOnFile)

	// Eventually write file
	if !write || reflect.DeepEqual(serversOnFile, allServers) {
		return allServers, nil
	}
	return allServers, s.FlushToFile(allServers)
}

func (s *storage) FlushToFile(servers models.AllServers) error {
	file, err := s.os.OpenFile(jsonFilepath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(servers); err != nil {
		_ = file.Close()
		return fmt.Errorf("cannot write to file: %w", err)
	}
	return file.Close()
}
