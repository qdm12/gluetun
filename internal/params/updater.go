package params

import (
	"time"

	libparams "github.com/qdm12/golibs/params"
)

// GetUpdaterPeriod obtains the period to fetch the servers information when the tunnel is up.
// Set to 0 to disable.
func (r *reader) GetUpdaterPeriod() (period time.Duration, err error) {
	s, err := r.envParams.GetEnv("UPDATER_PERIOD", libparams.Default("0"))
	if err != nil {
		return 0, err
	}
	return time.ParseDuration(s)
}
