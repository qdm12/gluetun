package params

import (
	libparams "github.com/qdm12/golibs/params"
)

func (r *reader) GetControlServerPort() (port uint16, warning string, err error) {
	return r.env.ListeningPort("HTTP_CONTROL_SERVER_PORT", libparams.Default("8000"))
}

func (r *reader) GetControlServerLog() (enabled bool, err error) {
	return r.env.OnOff("HTTP_CONTROL_SERVER_LOG", libparams.Default("on"))
}
