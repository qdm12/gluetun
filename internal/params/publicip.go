package params

import (
	"time"

	libparams "github.com/qdm12/golibs/params"
)

// GetPublicIPPeriod obtains the period to fetch the IP address periodically.
// Set to 0 to disable.
func (r *reader) GetPublicIPPeriod() (period time.Duration, err error) {
	s, err := r.envParams.GetEnv("PUBLICIP_PERIOD", libparams.Default("12h"))
	if err != nil {
		return 0, err
	}
	return time.ParseDuration(s)
}
