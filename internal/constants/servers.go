package constants

import (
	_ "embed"
	"encoding/json"

	"github.com/qdm12/gluetun/internal/models"
)

//go:embed servers.json
var allServersBytes []byte       //nolint:gochecknoglobals
var allServers models.AllServers //nolint:gochecknoglobals

func init() { //nolint:gochecknoinits
	// error returned covered by unit test
	allServers, _ = parseAllServers(allServersBytes)
}

func parseAllServers(b []byte) (allServers models.AllServers, err error) {
	err = json.Unmarshal(b, &allServers)
	return allServers, err
}

func GetAllServers() (allServers models.AllServers) {
	if allServers.Version == 0 { // not parsed yet - for unit tests mostly
		allServers, _ = parseAllServers(allServersBytes)
	}
	return allServers
}
