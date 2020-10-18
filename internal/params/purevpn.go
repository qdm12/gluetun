package params

import (
	"github.com/qdm12/gluetun/internal/constants"
)

// GetPurevpnRegions obtains the regions (continents) for the PureVPN servers from the
// environment variable REGION
func (r *reader) GetPurevpnRegions() (regions []string, err error) {
	choices := append(constants.PurevpnRegionChoices(), "")
	return r.envParams.GetCSVInPossibilities("REGION", choices)
}

// GetPurevpnCountries obtains the countries for the PureVPN servers from the
// environment variable COUNTRY
func (r *reader) GetPurevpnCountries() (countries []string, err error) {
	choices := append(constants.PurevpnCountryChoices(), "")
	return r.envParams.GetCSVInPossibilities("COUNTRY", choices)
}

// GetPurevpnCities obtains the cities for the PureVPN servers from the
// environment variable CITY
func (r *reader) GetPurevpnCities() (cities []string, err error) {
	choices := append(constants.PurevpnCityChoices(), "")
	return r.envParams.GetCSVInPossibilities("CITY", choices)
}
