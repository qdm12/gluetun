package params

import (
	"fmt"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	libparams "github.com/qdm12/golibs/params"
)

// GetWindscribeRegions obtains the regions for the Windscribe servers from the
// environment variable REGION.
func (r *reader) GetWindscribeRegions() (regions []string, err error) {
	return r.envParams.GetCSVInPossibilities("REGION", constants.WindscribeRegionChoices())
}

// GetWindscribeCities obtains the cities for the Windscribe servers from the
// environment variable CITY.
func (r *reader) GetWindscribeCities() (cities []string, err error) {
	return r.envParams.GetCSVInPossibilities("CITY", constants.WindscribeCityChoices())
}

// GetWindscribeHostnames obtains the hostnames for the Windscribe servers from the
// environment variable SERVER_HOSTNAME.
func (r *reader) GetWindscribeHostnames() (hostnames []string, err error) {
	return r.envParams.GetCSVInPossibilities("SERVER_HOSTNAME",
		constants.WindscribeHostnameChoices(),
		libparams.RetroKeys([]string{"HOSTNAME"}, r.onRetroActive),
	)
}

// GetWindscribePort obtains the port to reach the Windscribe server on from the
// environment variable PORT.
//nolint:gomnd
func (r *reader) GetWindscribePort(protocol models.NetworkProtocol) (port uint16, err error) {
	n, err := r.envParams.GetEnvIntRange("PORT", 0, 65535, libparams.Default("0"))
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
