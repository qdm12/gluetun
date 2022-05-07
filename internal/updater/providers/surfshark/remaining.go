package surfshark

import (
	"github.com/qdm12/gluetun/internal/provider/surfshark/servers"
)

// getRemainingServers finds extra servers not found in the API or in the ZIP file.
func getRemainingServers(hts hostToServer) {
	locationData := servers.LocationData()
	hostnameToLocationLeft := hostToLocation(locationData)
	for _, hostnameDone := range hts.toHostsSlice() {
		delete(hostnameToLocationLeft, hostnameDone)
	}

	for hostname, locationData := range hostnameToLocationLeft {
		// we assume the server supports TCP and UDP
		const tcp, udp = true, true
		hts.add(hostname, locationData.Region, locationData.Country,
			locationData.City, locationData.RetroLoc, tcp, udp)
	}
}
