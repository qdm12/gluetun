package portforward

import (
	"context"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/portforward/state"
)

func (l *Loop) GetStatus() (status models.LoopStatus) {
	return l.statusManager.GetStatus()
}

type StartData = state.StartData

func (l *Loop) Start(ctx context.Context, data StartData) (
	outcome string, err error) {
	l.startMu.Lock()
	defer l.startMu.Unlock()
	l.state.SetStartData(data)
	return l.statusManager.ApplyStatus(ctx, constants.Running)
}

func (l *Loop) Stop(ctx context.Context) (outcome string, err error) {
	return l.statusManager.ApplyStatus(ctx, constants.Stopped)
}
