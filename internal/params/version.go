package params

import (
	"github.com/qdm12/golibs/params"
)

func (p *paramsReader) GetVersion() string {
	version, _ := p.envParams.GetEnv("VERSION", params.Default("?"))
	return version
}

func (p *paramsReader) GetBuildDate() string {
	buildDate, _ := p.envParams.GetEnv("BUILD_DATE", params.Default("?"))
	return buildDate
}

func (p *paramsReader) GetVcsRef() string {
	buildDate, _ := p.envParams.GetEnv("VCS_REF", params.Default("?"))
	return buildDate
}
