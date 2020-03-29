package params

import (
	libparams "github.com/qdm12/golibs/params"
)

// GetUID obtains the user ID to use from the environment variable UID
func (p *paramsReader) GetUID() (uid int, err error) {
	return p.envParams.GetEnvIntRange("UID", 0, 65535, libparams.Default("1000"))
}

// GetGID obtains the group ID to use from the environment variable GID
func (p *paramsReader) GetGID() (gid int, err error) {
	return p.envParams.GetEnvIntRange("GID", 0, 65535, libparams.Default("1000"))
}
