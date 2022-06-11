package state

import (
	"context"
	"sync"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/models"
)

func New(statusApplier StatusApplier,
	settings settings.HTTPProxy) *State {
	return &State{
		statusApplier: statusApplier,
		settings:      settings,
	}
}

type State struct {
	statusApplier StatusApplier
	settings      settings.HTTPProxy
	settingsMu    sync.RWMutex
}

type StatusApplier interface {
	ApplyStatus(ctx context.Context, status models.LoopStatus) (
		outcome string, err error)
}
