package params

import (
	"strings"

	"github.com/qdm12/private-internet-access-docker/internal/constants"
	"github.com/qdm12/private-internet-access-docker/internal/models"
)

// GetSurfsharkRegion obtains the region for the Surfshark server from the
// environment variable REGION
func (r *reader) GetSurfsharkRegion() (region models.SurfsharkRegion, err error) {
	s, err := r.envParams.GetValueIfInside("REGION", constants.SurfsharkRegionChoices())
	return models.SurfsharkRegion(strings.ToLower(s)), err
}
