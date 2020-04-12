package params

import (
	"strings"

	libparams "github.com/qdm12/golibs/params"
	"github.com/qdm12/private-internet-access-docker/internal/constants"
	"github.com/qdm12/private-internet-access-docker/internal/models"
)

// GetMullvadCountry obtains the country for the Mullvad server from the
// environment variable COUNTRY
func (p *reader) GetMullvadCountry() (country models.MullvadCountry, err error) {
	choices := append(constants.MullvadCountryChoices(), "")
	s, err := p.envParams.GetValueIfInside("COUNTRY", choices)
	return models.MullvadCountry(strings.ToLower(s)), err
}

// GetMullvadCity obtains the city for the Mullvad server from the
// environment variable CITY
func (p *reader) GetMullvadCity() (country models.MullvadCity, err error) {
	choices := append(constants.MullvadCityChoices(), "")
	s, err := p.envParams.GetValueIfInside("CITY", choices)
	return models.MullvadCity(strings.ToLower(s)), err
}

// GetMullvadISP obtains the ISP for the Mullvad server from the
// environment variable ISP
func (p *reader) GetMullvadISP() (country models.MullvadProvider, err error) {
	choices := append(constants.MullvadProviderChoices(), "")
	s, err := p.envParams.GetValueIfInside("ISP", choices)
	return models.MullvadProvider(strings.ToLower(s)), err
}

// GetMullvadPort obtains the port to reach the Mullvad server on from the
// environment variable PORT
func (p *reader) GetMullvadPort() (port uint16, err error) {
	n, err := p.envParams.GetEnvIntRange("PORT", 0, 65535, libparams.Default("0"))
	return uint16(n), err
}
