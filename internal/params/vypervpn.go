package params

import (
	"github.com/qdm12/gluetun/internal/constants"
)

// GetVyprvpnRegion obtains the region for the Vyprvpn server from the
// environment variable REGION
func (r *reader) GetVyprvpnRegion() (region string, err error) {
	choices := append(constants.VyprvpnRegionChoices(), "")
	return r.envParams.GetValueIfInside("REGION", choices)
}
