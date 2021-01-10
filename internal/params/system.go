package params

import (
	libparams "github.com/qdm12/golibs/params"
)

// GetPUID obtains the user ID to use from the environment variable PUID
// with retro compatible variable UID.
func (r *reader) GetPUID() (ppuid int, err error) {
	return r.env.IntRange("PUID", 0, 65535,
		libparams.Default("1000"),
		libparams.RetroKeys([]string{"UID"}, r.onRetroActive))
}

// GetGID obtains the group ID to use from the environment variable PGID
// with retro compatible variable PGID.
func (r *reader) GetPGID() (pgid int, err error) {
	return r.env.IntRange("PGID", 0, 65535,
		libparams.Default("1000"),
		libparams.RetroKeys([]string{"GID"}, r.onRetroActive))
}

// GetTZ obtains the timezone from the environment variable TZ.
func (r *reader) GetTimezone() (timezone string, err error) {
	return r.env.Get("TZ")
}
