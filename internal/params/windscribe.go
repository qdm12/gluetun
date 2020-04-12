package params

import (
	"fmt"
	"strings"

	libparams "github.com/qdm12/golibs/params"
	"github.com/qdm12/private-internet-access-docker/internal/constants"
	"github.com/qdm12/private-internet-access-docker/internal/models"
)

// GetWindscribeRegion obtains the region for the Windscribe server from the
// environment variable REGION
func (p *reader) GetWindscribeRegion() (country models.WindscribeRegion, err error) {
	s, err := p.envParams.GetValueIfInside("REGION", constants.WindscribeRegionChoices())
	return models.WindscribeRegion(strings.ToLower(s)), err
}

// GetMullvadPort obtains the port to reach the Mullvad server on from the
// environment variable PORT
func (p *reader) GetWindscribePort(protocol models.NetworkProtocol) (port uint16, err error) {
	n, err := p.envParams.GetEnvIntRange("PORT", 0, 65535, libparams.Default("0"))
	if err != nil {
		return 0, err
	}
	if n == 0 {
		return 0, nil
	}
	switch protocol {
	case constants.TCP:
		switch n {
		case 21, 22, 80, 123, 143, 443, 587, 1194, 3306, 8080, 54783:
		default:
			return 0, fmt.Errorf("port %d is not valid for protocol %s", n, protocol)
		}
	case constants.UDP:
		switch n {
		case 53, 80, 123, 443, 1194, 54783:
		default:
			return 0, fmt.Errorf("port %d is not valid for protocol %s", n, protocol)
		}
	}
	return uint16(n), nil
}
