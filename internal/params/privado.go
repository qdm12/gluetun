package params

import (
	"github.com/qdm12/gluetun/internal/constants"
	libparams "github.com/qdm12/golibs/params"
)

// GetPrivadoHostnames obtains the hostnames for the Privado server from the
// environment variable SERVER_HOSTNAME.
func (r *reader) GetPrivadoHostnames() (hosts []string, err error) {
	return r.env.CSVInside("SERVER_HOSTNAME",
		constants.PrivadoHostnameChoices(),
		libparams.RetroKeys([]string{"HOSTNAME"}, r.onRetroActive))
}
