package params

import (
	"github.com/qdm12/gluetun/internal/constants"
	libparams "github.com/qdm12/golibs/params"
)

// GetMullvadCountry obtains the country for the Mullvad server from the
// environment variable COUNTRY
func (r *reader) GetMullvadCountry() (country string, err error) {
	choices := append(constants.MullvadCountryChoices(), "")
	return r.envParams.GetValueIfInside("COUNTRY", choices)
}

// GetMullvadCity obtains the city for the Mullvad server from the
// environment variable CITY
func (r *reader) GetMullvadCity() (country string, err error) {
	choices := append(constants.MullvadCityChoices(), "")
	return r.envParams.GetValueIfInside("CITY", choices)
}

// GetMullvadISP obtains the ISP for the Mullvad server from the
// environment variable ISP
func (r *reader) GetMullvadISP() (isp string, err error) {
	choices := append(constants.MullvadISPChoices(), "")
	return r.envParams.GetValueIfInside("ISP", choices)
}

// GetMullvadPort obtains the port to reach the Mullvad server on from the
// environment variable PORT
func (r *reader) GetMullvadPort() (port uint16, err error) {
	n, err := r.envParams.GetEnvIntRange("PORT", 0, 65535, libparams.Default("0"))
	return uint16(n), err
}
