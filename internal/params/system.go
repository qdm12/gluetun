package params

import (
	libparams "github.com/qdm12/golibs/params"
)

// GetUID obtains the user ID to use from the environment variable UID.
func (r *reader) GetUID() (uid int, err error) {
	return r.envParams.GetEnvIntRange("UID", 0, 65535, libparams.Default("1000"))
}

// GetGID obtains the group ID to use from the environment variable GID.
func (r *reader) GetGID() (gid int, err error) {
	return r.envParams.GetEnvIntRange("GID", 0, 65535, libparams.Default("1000"))
}

// GetTZ obtains the timezone from the environment variable TZ.
func (r *reader) GetTimezone() (timezone string, err error) {
	return r.envParams.GetEnv("TZ")
}
