package params

import (
	libparams "github.com/qdm12/golibs/params"
)

func (r *reader) GetControlServerPort() (port uint16, err error) {
	n, err := r.envParams.GetEnvIntRange("HTTP_CONTROL_SERVER_PORT", 1, 65535, libparams.Default("8000"))
	if err != nil {
		return 0, err
	}
	return uint16(n), nil
}

func (r *reader) GetControlServerLog() (enabled bool, err error) {
	return r.envParams.GetOnOff("HTTP_CONTROL_SERVER_LOG", libparams.Default("on"))
}
