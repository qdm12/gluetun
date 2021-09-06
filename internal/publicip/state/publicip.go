package state

import (
	"github.com/qdm12/gluetun/internal/publicip/models"
)

type DataGetSetter interface {
	GetData() (data models.IPInfoData)
	SetData(data models.IPInfoData)
}

func (s *State) GetData() (data models.IPInfoData) {
	s.ipDataMu.RLock()
	defer s.ipDataMu.RUnlock()
	return s.ipData.Copy()
}

func (s *State) SetData(data models.IPInfoData) {
	s.ipDataMu.Lock()
	defer s.ipDataMu.Unlock()
	s.ipData = data.Copy()
}
