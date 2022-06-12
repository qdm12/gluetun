package publicip

import "github.com/qdm12/gluetun/internal/publicip/ipinfo"

func (l *Loop) GetData() (data ipinfo.Response) {
	return l.state.GetData()
}

func (l *Loop) SetData(data ipinfo.Response) {
	l.state.SetData(data)
}
