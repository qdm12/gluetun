package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"reflect"

	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/golibs/os"
)

var (
	ErrCannotReadFile  = errors.New("cannot read servers from file")
	ErrCannotWriteFile = errors.New("cannot write servers to file")
)

func countServers(allServers models.AllServers) int {
	return len(allServers.Cyberghost.Servers) +
		len(allServers.Fastestvpn.Servers) +
		len(allServers.HideMyAss.Servers) +
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
		len(allServers.Vyprvpn.Servers) +
		len(allServers.Windscribe.Servers)
}

func (s *storage) SyncServers(hardcodedServers models.AllServers) (
	allServers models.AllServers, err error) {
	serversOnFile, err := s.readFromFile(s.filepath)
	if err != nil {
		return allServers, fmt.Errorf("%w: %s", ErrCannotReadFile, err)
	}

	hardcodedCount := countServers(hardcodedServers)
	countOnFile := countServers(serversOnFile)

	if countOnFile == 0 {
		s.logger.Info("creating %s with %d hardcoded servers", s.filepath, hardcodedCount)
		allServers = hardcodedServers
	} else {
		s.logger.Info(
			"merging by most recent %d hardcoded servers and %d servers read from %s",
			hardcodedCount, countOnFile, s.filepath)
		allServers = s.mergeServers(hardcodedServers, serversOnFile)
	}

	// Eventually write file
	if s.filepath == "" || reflect.DeepEqual(serversOnFile, allServers) {
		return allServers, nil
	}

	if err := s.FlushToFile(allServers); err != nil {
		return allServers, fmt.Errorf("%w: %s", ErrCannotWriteFile, err)
	}
	return allServers, nil
}

func (s *storage) readFromFile(filepath string) (servers models.AllServers, err error) {
	file, err := s.os.OpenFile(filepath, os.O_RDONLY, 0)
	if os.IsNotExist(err) {
		return servers, nil
	} else if err != nil {
		return servers, err
	}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&servers); err != nil {
		_ = file.Close()
		if errors.Is(err, io.EOF) {
			return servers, nil
		}
		return servers, err
	}
	return servers, file.Close()
}

func (s *storage) FlushToFile(servers models.AllServers) error {
	dirPath := filepath.Dir(s.filepath)
	if err := s.os.MkdirAll(dirPath, 0644); err != nil {
		return err
	}

	file, err := s.os.OpenFile(s.filepath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(servers); err != nil {
		_ = file.Close()
		return err
	}
	return file.Close()
}
