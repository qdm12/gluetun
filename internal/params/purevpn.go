package params

import (
	"github.com/qdm12/gluetun/internal/constants"
)

// GetPurevpnRegions obtains the regions (continents) for the PureVPN servers from the
// environment variable REGION.
func (r *reader) GetPurevpnRegions() (regions []string, err error) {
	return r.envParams.GetCSVInPossibilities("REGION", constants.PurevpnRegionChoices())
}

// GetPurevpnCountries obtains the countries for the PureVPN servers from the
// environment variable COUNTRY.
func (r *reader) GetPurevpnCountries() (countries []string, err error) {
	return r.envParams.GetCSVInPossibilities("COUNTRY", constants.PurevpnCountryChoices())
}

// GetPurevpnCities obtains the cities for the PureVPN servers from the
// environment variable CITY.
func (r *reader) GetPurevpnCities() (cities []string, err error) {
	return r.envParams.GetCSVInPossibilities("CITY", constants.PurevpnCityChoices())
}
