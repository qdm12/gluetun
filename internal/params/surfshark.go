package params

import (
	"github.com/qdm12/private-internet-access-docker/internal/constants"
)

// GetSurfsharkRegion obtains the region for the Surfshark server from the
// environment variable REGION
func (r *reader) GetSurfsharkRegion() (region string, err error) {
	choices := append(constants.SurfsharkRegionChoices(), "")
	return r.envParams.GetValueIfInside("REGION", choices)
}
