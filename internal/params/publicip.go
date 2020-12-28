package params

import (
	"time"

	"github.com/qdm12/gluetun/internal/models"
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

// GetPublicIPFilepath obtains the public IP filepath
// from the environment variable PUBLICIP_FILE with retro-compatible
// environment variable IP_STATUS_FILE.
func (r *reader) GetPublicIPFilepath() (filepath models.Filepath, err error) {
	filepathStr, err := r.envParams.GetPath("PUBLICIP_FILE",
		libparams.RetroKeys([]string{"IP_STATUS_FILE"}, r.onRetroActive),
		libparams.Default("/tmp/gluetun/ip"), libparams.CaseSensitiveValue())
	return models.Filepath(filepathStr), err
}
