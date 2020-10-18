package params

import (
	"github.com/qdm12/gluetun/internal/constants"
)

// GetSurfsharkRegions obtains the regions for the Surfshark servers from the
// environment variable REGION
func (r *reader) GetSurfsharkRegions() (regions []string, err error) {
	choices := append(constants.SurfsharkRegionChoices(), "")
	return r.envParams.GetCSVInPossibilities("REGION", choices)
}
