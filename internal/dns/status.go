package dns

import (
	"context"

	"github.com/qdm12/gluetun/internal/models"
)

func (l *looper) GetStatus() (status models.LoopStatus) { return l.state.GetStatus() }

type StatusApplier interface {
	ApplyStatus(ctx context.Context, status models.LoopStatus) (
		outcome string, err error)
}

func (l *looper) ApplyStatus(ctx context.Context, status models.LoopStatus) (
	outcome string, err error) {
	return l.state.ApplyStatus(ctx, status)
}
