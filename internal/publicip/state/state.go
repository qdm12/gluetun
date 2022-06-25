package state

import (
	"context"
	"sync"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/models"
)

func New(statusApplier StatusApplier,
	settings settings.PublicIP,
	updateTicker chan<- struct{}) *State {
	return &State{
		statusApplier: statusApplier,
		settings:      settings,
		updateTicker:  updateTicker,
	}
}

type State struct {
	statusApplier StatusApplier

	settings   settings.PublicIP
	settingsMu sync.RWMutex

	ipData   models.PublicIP
	ipDataMu sync.RWMutex

	updateTicker chan<- struct{}
}

type StatusApplier interface {
	ApplyStatus(ctx context.Context, status models.LoopStatus) (
		outcome string, err error)
}
