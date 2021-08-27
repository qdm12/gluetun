package storage

import (
	"encoding/json"
	"errors"
	"io"
	"os"

	"github.com/qdm12/gluetun/internal/models"
)

func readFromFile(filepath string) (servers models.AllServers, err error) {
	file, err := os.Open(filepath)
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
