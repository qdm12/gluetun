package storage

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/qdm12/gluetun/internal/models"
)

var (
	ErrCannotReadFile  = errors.New("cannot read servers from file")
	ErrCannotWriteFile = errors.New("cannot write servers to file")
)

func countServers(allServers models.AllServers) int {
	return len(allServers.Cyberghost.Servers) +
		len(allServers.Fastestvpn.Servers) +
		len(allServers.HideMyAss.Servers) +
		len(allServers.Ipvanish.Servers) +
		len(allServers.Ivpn.Servers) +
		len(allServers.Mullvad.Servers) +
		len(allServers.Nordvpn.Servers) +
		len(allServers.Privado.Servers) +
		len(allServers.Pia.Servers) +
		len(allServers.Privatevpn.Servers) +
		len(allServers.Protonvpn.Servers) +
		len(allServers.Purevpn.Servers) +
		len(allServers.Surfshark.Servers) +
		len(allServers.Torguard.Servers) +
		len(allServers.VPNUnlimited.Servers) +
		len(allServers.Vyprvpn.Servers) +
		len(allServers.Wevpn.Servers) +
		len(allServers.Windscribe.Servers)
}

func (s *Storage) SyncServers() (err error) {
	serversOnFile, err := readFromFile(s.filepath)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrCannotReadFile, err)
	}

	hardcodedCount := countServers(s.hardcodedServers)
	countOnFile := countServers(serversOnFile)

	if countOnFile == 0 {
		s.logger.Info(fmt.Sprintf(
			"creating %s with %d hardcoded servers",
			s.filepath, hardcodedCount))
		s.mergedServers = s.hardcodedServers
	} else {
		s.logger.Info(fmt.Sprintf(
			"merging by most recent %d hardcoded servers and %d servers read from %s",
			hardcodedCount, countOnFile, s.filepath))

		s.mergedServers = s.mergeServers(s.hardcodedServers, serversOnFile)
	}

	// Eventually write file
	if s.filepath == "" || reflect.DeepEqual(serversOnFile, s.mergedServers) {
		return nil
	}

	if err := flushToFile(s.filepath, s.mergedServers); err != nil {
		return fmt.Errorf("%w: %s", ErrCannotWriteFile, err)
	}
	return nil
}
