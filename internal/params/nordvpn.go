package params

import (
	"fmt"
	"strconv"

	"github.com/qdm12/gluetun/internal/constants"
)

// GetNordvpnRegions obtains the regions (countries) for the NordVPN server from the
// environment variable REGION.
func (r *reader) GetNordvpnRegions() (regions []string, err error) {
	return r.envParams.GetCSVInPossibilities("REGION", constants.NordvpnRegionChoices())
}

// GetNordvpnRegion obtains the server numbers (optional) for the NordVPN servers from the
// environment variable SERVER_NUMBER.
func (r *reader) GetNordvpnNumbers() (numbers []uint16, err error) {
	possibilities := make([]string, 65537)
	for i := range possibilities {
		possibilities[i] = fmt.Sprintf("%d", i)
	}
	possibilities[65536] = ""
	values, err := r.envParams.GetCSVInPossibilities("SERVER_NUMBER", possibilities)
	if err != nil {
		return nil, err
	}
	numbers = make([]uint16, len(values))
	for i := range values {
		n, err := strconv.Atoi(values[i])
		if err != nil {
			return nil, err
		}
		numbers[i] = uint16(n)
	}
	return numbers, nil
}
