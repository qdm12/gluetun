package params

import (
	"github.com/qdm12/private-internet-access-docker/internal/constants"
)

// GetNordvpnRegion obtains the region (server name) for the NordVPN server from the
// environment variable REGION
func (r *reader) GetNordvpnRegion() (region string, err error) {
	return r.envParams.GetValueIfInside("REGION", constants.NordvpnRegionChoices())
}
