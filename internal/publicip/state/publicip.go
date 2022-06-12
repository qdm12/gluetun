package state

import (
	"github.com/qdm12/gluetun/internal/publicip/ipinfo"
)

func (s *State) GetData() (data ipinfo.Response) {
	s.ipDataMu.RLock()
	defer s.ipDataMu.RUnlock()
	return s.ipData.Copy()
}

func (s *State) SetData(data ipinfo.Response) {
	s.ipDataMu.Lock()
	defer s.ipDataMu.Unlock()
	s.ipData = data.Copy()
}
