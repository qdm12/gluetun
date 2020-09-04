package params

import (
	"github.com/qdm12/gluetun/internal/constants"
)

// GetPurevpnRegion obtains the region (continent) for the PureVPN server from the
// environment variable REGION
func (r *reader) GetPurevpnRegion() (region string, err error) {
	choices := append(constants.PurevpnRegionChoices(), "")
	return r.envParams.GetValueIfInside("REGION", choices)
}

// GetPurevpnCountry obtains the country for the PureVPN server from the
// environment variable COUNTRY
func (r *reader) GetPurevpnCountry() (country string, err error) {
	choices := append(constants.PurevpnCountryChoices(), "")
	return r.envParams.GetValueIfInside("COUNTRY", choices)
}

// GetPurevpnCity obtains the city for the PureVPN server from the
// environment variable CITY
func (r *reader) GetPurevpnCity() (city string, err error) {
	choices := append(constants.PurevpnCityChoices(), "")
	return r.envParams.GetValueIfInside("CITY", choices)
}
