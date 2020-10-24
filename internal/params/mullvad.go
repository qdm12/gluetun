package params

import (
	"github.com/qdm12/gluetun/internal/constants"
	libparams "github.com/qdm12/golibs/params"
)

// GetMullvadCountries obtains the countries for the Mullvad servers from the
// environment variable COUNTRY.
func (r *reader) GetMullvadCountries() (countries []string, err error) {
	return r.envParams.GetCSVInPossibilities("COUNTRY", constants.MullvadCountryChoices())
}

// GetMullvadCity obtains the cities for the Mullvad servers from the
// environment variable CITY.
func (r *reader) GetMullvadCities() (cities []string, err error) {
	return r.envParams.GetCSVInPossibilities("CITY", constants.MullvadCityChoices())
}

// GetMullvadISPs obtains the ISPs for the Mullvad servers from the
// environment variable ISP.
func (r *reader) GetMullvadISPs() (isps []string, err error) {
	return r.envParams.GetCSVInPossibilities("ISP", constants.MullvadISPChoices())
}

// GetMullvadPort obtains the port to reach the Mullvad server on from the
// environment variable PORT.
func (r *reader) GetMullvadPort() (port uint16, err error) {
	n, err := r.envParams.GetEnvIntRange("PORT", 0, 65535, libparams.Default("0"))
	return uint16(n), err
}

// GetMullvadOwned obtains if the server should be owned by Mullvad or not from the
// environment variable OWNED.
func (r *reader) GetMullvadOwned() (owned bool, err error) {
	return r.envParams.GetYesNo("OWNED", libparams.Default("no"))
}
