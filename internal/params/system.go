package params

import (
	libparams "github.com/qdm12/golibs/params"
	"github.com/qdm12/private-internet-access-docker/internal/models"
)

// GetUID obtains the user ID to use from the environment variable UID
func (p *reader) GetUID() (uid int, err error) {
	return p.envParams.GetEnvIntRange("UID", 0, 65535, libparams.Default("1000"))
}

// GetGID obtains the group ID to use from the environment variable GID
func (p *reader) GetGID() (gid int, err error) {
	return p.envParams.GetEnvIntRange("GID", 0, 65535, libparams.Default("1000"))
}

// GetTZ obtains the timezone from the environment variable TZ
func (p *reader) GetTimezone() (timezone string, err error) {
	return p.envParams.GetEnv("TZ")
}

// GetIPStatusFilepath obtains the IP status file path
// from the environment variable IP_STATUS_FILE
func (p *reader) GetIPStatusFilepath() (filepath models.Filepath, err error) {
	filepathStr, err := p.envParams.GetPath("IP_STATUS_FILE", libparams.Default("/ip"), libparams.CaseSensitiveValue())
	return models.Filepath(filepathStr), err
}
