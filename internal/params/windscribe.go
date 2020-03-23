package params

import (
	"strings"

	"github.com/qdm12/private-internet-access-docker/internal/constants"
	"github.com/qdm12/private-internet-access-docker/internal/models"
)

// GetWindscribeRegion obtains the region for the Windscribe server from the
// environment variable REGION
func (p *paramsReader) GetWindscribeRegion() (country models.WindscribeRegion, err error) {
	choices := append(constants.WindscribeRegionChoices())
	s, err := p.envParams.GetValueIfInside("REGION", choices)
	return models.WindscribeRegion(strings.ToLower(s)), err
}
