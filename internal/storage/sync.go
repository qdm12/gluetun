package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	"github.com/qdm12/gluetun/internal/models"
)

const (
	jsonFilepath = "/gluetun/servers.json"
)

func countServers(allServers models.AllServers) int {
	return len(allServers.Cyberghost.Servers) +
		len(allServers.Mullvad.Servers) +
		len(allServers.Nordvpn.Servers) +
		len(allServers.Pia.Servers) +
		len(allServers.PiaOld.Servers) +
		len(allServers.Purevpn.Servers) +
		len(allServers.Surfshark.Servers) +
		len(allServers.Vyprvpn.Servers) +
		len(allServers.Windscribe.Servers)
}

func (s *storage) SyncServers(hardcodedServers models.AllServers, write bool) (allServers models.AllServers, err error) {
	// Eventually read file
	var serversOnFile models.AllServers
	_, err = s.osStat(jsonFilepath)
	if err == nil {
		serversOnFile, err = s.readFromFile()
		if err != nil {
			return allServers, err
		}
	} else if !os.IsNotExist(err) {
		return allServers, err
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

func (s *storage) readFromFile() (servers models.AllServers, err error) {
	bytes, err := s.readFile(jsonFilepath)
	if err != nil {
		return servers, err
	}
	if err := json.Unmarshal(bytes, &servers); err != nil {
		return servers, err
	}
	return servers, nil
}

func (s *storage) FlushToFile(servers models.AllServers) error {
	bytes, err := json.MarshalIndent(servers, "", "  ")
	if err != nil {
		return fmt.Errorf("cannot write to file: %w", err)
	}
	if err := s.writeFile(jsonFilepath, bytes, 0644); err != nil {
		return err
	}
	return nil
}
