package params

import (
	"github.com/qdm12/gluetun/internal/constants"
	libparams "github.com/qdm12/golibs/params"
)

// GetNordvpnRegion obtains the region (country) for the NordVPN server from the
// environment variable REGION
func (r *reader) GetNordvpnRegion() (region string, err error) {
	choices := append(constants.NordvpnRegionChoices(), "")
	return r.envParams.GetValueIfInside("REGION", choices)
}

// GetNordvpnRegion obtains the server number (optional) for the NordVPN server from the
// environment variable SERVER_NUMBER
func (r *reader) GetNordvpnNumber() (number uint16, err error) {
	n, err := r.envParams.GetEnvIntRange("SERVER_NUMBER", 0, 65535, libparams.Default("0"))
	if err != nil {
		return 0, err
	}
	return uint16(n), nil
}
