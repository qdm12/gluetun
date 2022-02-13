package surfshark

import (
	"errors"
	"fmt"

	"github.com/qdm12/gluetun/internal/models"
)

var (
	errHostnameNotFound = errors.New("hostname not found in hostname to location mapping")
)

func getHostInformation(host string, hostnameToLocation map[string]models.SurfsharkLocationData) (
	data models.SurfsharkLocationData, err error) {
	locationData, ok := hostnameToLocation[host]
	if !ok {
		return locationData, fmt.Errorf("%w: %s", errHostnameNotFound, host)
	}

	return locationData, nil
}

func hostToLocation(locationData []models.SurfsharkLocationData) (
	hostToLocation map[string]models.SurfsharkLocationData) {
	hostToLocation = make(map[string]models.SurfsharkLocationData, len(locationData))
	for _, data := range locationData {
		hostToLocation[data.Hostname] = data
	}
	return hostToLocation
}
