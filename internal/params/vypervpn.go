package params

import (
	"github.com/qdm12/private-internet-access-docker/internal/constants"
)

// GetVyprvpnRegion obtains the region for the Vyprvpn server from the
// environment variable REGION
func (r *reader) GetVyprvpnRegion() (region string, err error) {
	return r.envParams.GetValueIfInside("REGION", constants.VyprvpnRegionChoices())
}
