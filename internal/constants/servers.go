package constants

import (
	"embed"
	"encoding/json"
	"sync"

	"github.com/qdm12/gluetun/internal/models"
)

//go:embed servers.json
var allServersEmbedFS embed.FS   //nolint:gochecknoglobals
var allServers models.AllServers //nolint:gochecknoglobals
var parseOnce sync.Once          //nolint:gochecknoglobals

func init() { //nolint:gochecknoinits
	// error returned covered by unit test
	parseOnce.Do(func() { allServers, _ = parseAllServers() })
}

func parseAllServers() (allServers models.AllServers, err error) {
	f, err := allServersEmbedFS.Open("servers.json")
	if err != nil {
		return allServers, err
	}
	decoder := json.NewDecoder(f)
	err = decoder.Decode(&allServers)
	return allServers, err
}

func GetAllServers() models.AllServers {
	parseOnce.Do(func() { allServers, _ = parseAllServers() }) // init did not execute, used in tests
	return allServers
}
