package state

import (
	"context"
	"sync"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/models"
)

func New(statusApplier StatusApplier, vpn settings.VPN) *State {
	return &State{
		statusApplier: statusApplier,
		vpn:           vpn,
	}
}

type State struct {
	statusApplier StatusApplier

	vpn        settings.VPN
	settingsMu sync.RWMutex
}

type StatusApplier interface {
	ApplyStatus(ctx context.Context, status models.LoopStatus) (
		outcome string, err error)
}
