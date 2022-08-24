package updater

import (
	"github.com/qdm12/gluetun/internal/provider/surfshark/servers"
)

// getRemainingServers finds extra servers not found in the API or in the ZIP file.
func getRemainingServers(hts hostToServers) {
	locationData := servers.LocationData()
	hostnameToLocationLeft := hostToLocation(locationData)
	for _, hostnameDone := range hts.toHostsSlice() {
		delete(hostnameToLocationLeft, hostnameDone)
	}

	for hostname, locationData := range hostnameToLocationLeft {
		// we assume the OpenVPN server supports both TCP and UDP
		const tcp, udp = true, true
		hts.addOpenVPN(hostname, locationData.Region, locationData.Country,
			locationData.City, locationData.RetroLoc, tcp, udp)
	}
}
