package params

import (
	libparams "github.com/qdm12/golibs/params"
)

func (p *paramsReader) GetVersion() string {
	version, _ := p.envParams.GetEnv("VERSION", libparams.Default("?"), libparams.CaseSensitiveValue())
	return version
}

func (p *paramsReader) GetBuildDate() string {
	buildDate, _ := p.envParams.GetEnv("BUILD_DATE", libparams.Default("?"), libparams.CaseSensitiveValue())
	return buildDate
}

func (p *paramsReader) GetVcsRef() string {
	buildDate, _ := p.envParams.GetEnv("VCS_REF", libparams.Default("?"), libparams.CaseSensitiveValue())
	return buildDate
}
