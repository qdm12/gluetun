package state

import (
	"sync"

	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/loopstate"
	"github.com/qdm12/gluetun/internal/publicip/models"
)

var _ Manager = (*State)(nil)

type Manager interface {
	SettingsGetSetter
	DataGetSetter
}

func New(statusApplier loopstate.Applier,
	settings configuration.PublicIP,
	updateTicker chan<- struct{}) *State {
	return &State{
		statusApplier: statusApplier,
		settings:      settings,
		updateTicker:  updateTicker,
	}
}

type State struct {
	statusApplier loopstate.Applier

	settings   configuration.PublicIP
	settingsMu sync.RWMutex

	ipData   models.IPInfoData
	ipDataMu sync.RWMutex

	updateTicker chan<- struct{}
}
