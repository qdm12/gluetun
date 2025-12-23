package storage

import (
	"embed"
	"encoding/json"

	"github.com/qdm12/gluetun/internal/models"
)

//go:embed servers.json
var allServersEmbedFS embed.FS

func parseHardcodedServers() (allServers models.AllServers, err error) {
	f, err := allServersEmbedFS.Open("servers.json")
	if err != nil {
		return allServers, err
	}
	decoder := json.NewDecoder(f)
	err = decoder.Decode(&allServers)
	return allServers, err
}
