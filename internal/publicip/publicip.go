package publicip

import (
	"github.com/qdm12/gluetun/internal/models"
)

func (l *Loop) GetData() (data models.PublicIP) {
	return l.state.GetData()
}

func (l *Loop) SetData(data models.PublicIP) {
	l.state.SetData(data)
}
