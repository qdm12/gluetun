package params

import (
	"github.com/qdm12/gluetun/internal/constants"
)

// GetPrivadoHostnames obtains the hostnames for the Privado server from the
// environment variable HOSTNAME.
func (r *reader) GetPrivadoHostnames() (hosts []string, err error) {
	return r.envParams.GetCSVInPossibilities("HOSTNAME", constants.PrivadoHostnameChoices())
}
