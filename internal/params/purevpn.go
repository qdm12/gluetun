package params

import (
	"github.com/qdm12/private-internet-access-docker/internal/constants"
)

// GetPurevpnRegion obtains the region (continent) for the PureVPN server from the
// environment variable REGION
func (r *reader) GetPurevpnRegion() (region string, err error) {
	return r.envParams.GetValueIfInside("REGION", constants.PurevpnRegionChoices())
}

// GetPurevpnCountry obtains the country for the PureVPN server from the
// environment variable COUNTRY
func (r *reader) GetPurevpnCountry() (country string, err error) {
	return r.envParams.GetValueIfInside("COUNTRY", constants.PurevpnCountryChoices())
}

// GetPurevpnCity obtains the city for the PureVPN server from the
// environment variable CITY
func (r *reader) GetPurevpnCity() (city string, err error) {
	return r.envParams.GetValueIfInside("CITY", constants.PurevpnCityChoices())
}
