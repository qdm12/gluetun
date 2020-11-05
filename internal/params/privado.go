package params

import (
	"fmt"
	"strconv"

	"github.com/qdm12/gluetun/internal/constants"
)

// GetPrivadoCities obtains the cities for the Privado server from the
// environment variable CITY.
func (r *reader) GetPrivadoCities() (regions []string, err error) {
	return r.envParams.GetCSVInPossibilities("CITY", constants.PrivadoCityChoices())
}

// GetPrivadoNumbers obtains the server numbers (optional) for the Privado servers from the
// environment variable SERVER_NUMBER.
func (r *reader) GetPrivadoNumbers() (numbers []uint16, err error) {
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
