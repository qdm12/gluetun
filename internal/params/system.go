package params

import (
	"github.com/qdm12/gluetun/internal/models"
	libparams "github.com/qdm12/golibs/params"
)

// GetUID obtains the user ID to use from the environment variable UID
func (r *reader) GetUID() (uid int, err error) {
	return r.envParams.GetEnvIntRange("UID", 0, 65535, libparams.Default("1000"))
}

// GetGID obtains the group ID to use from the environment variable GID
func (r *reader) GetGID() (gid int, err error) {
	return r.envParams.GetEnvIntRange("GID", 0, 65535, libparams.Default("1000"))
}

// GetTZ obtains the timezone from the environment variable TZ
func (r *reader) GetTimezone() (timezone string, err error) {
	return r.envParams.GetEnv("TZ")
}

// GetIPStatusFilepath obtains the IP status file path
// from the environment variable IP_STATUS_FILE
func (r *reader) GetIPStatusFilepath() (filepath models.Filepath, err error) {
	filepathStr, err := r.envParams.GetPath("IP_STATUS_FILE", libparams.Default("/tmp/gluetun/ip"), libparams.CaseSensitiveValue())
	return models.Filepath(filepathStr), err
}
