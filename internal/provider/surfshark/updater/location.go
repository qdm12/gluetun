package updater

import (
	"errors"
	"fmt"

	"github.com/qdm12/gluetun/internal/provider/surfshark/servers"
)

var (
	errHostnameNotFound = errors.New("hostname not found in hostname to location mapping")
)

func getHostInformation(host string, hostnameToLocation map[string]servers.ServerLocation) (
	data servers.ServerLocation, err error) {
	locationData, ok := hostnameToLocation[host]
	if !ok {
		return locationData, fmt.Errorf("%w: %s", errHostnameNotFound, host)
	}

	return locationData, nil
}

func hostToLocation(locationData []servers.ServerLocation) (
	hostToLocation map[string]servers.ServerLocation) {
	hostToLocation = make(map[string]servers.ServerLocation, len(locationData))
	for _, data := range locationData {
		hostToLocation[data.Hostname] = data
	}
	return hostToLocation
}
