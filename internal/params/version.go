package params

import (
	libparams "github.com/qdm12/golibs/params"
)

func (r *reader) GetVersion() string {
	version, _ := r.envParams.GetEnv("VERSION", libparams.Default("?"), libparams.CaseSensitiveValue())
	return version
}

func (r *reader) GetBuildDate() string {
	buildDate, _ := r.envParams.GetEnv("BUILD_DATE", libparams.Default("?"), libparams.CaseSensitiveValue())
	return buildDate
}

func (r *reader) GetVcsRef() string {
	buildDate, _ := r.envParams.GetEnv("VCS_REF", libparams.Default("?"), libparams.CaseSensitiveValue())
	return buildDate
}
