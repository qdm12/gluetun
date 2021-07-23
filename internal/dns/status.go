package dns

import (
	"context"

	"github.com/qdm12/gluetun/internal/models"
)

type StatusGetterApplier interface {
	StatusGetter
	StatusApplier
}

type StatusGetter interface {
	GetStatus() (status models.LoopStatus)
}

func (l *Loop) GetStatus() (status models.LoopStatus) { return l.state.GetStatus() }

type StatusApplier interface {
	ApplyStatus(ctx context.Context, status models.LoopStatus) (
		outcome string, err error)
}

func (l *Loop) ApplyStatus(ctx context.Context, status models.LoopStatus) (
	outcome string, err error) {
	return l.state.ApplyStatus(ctx, status)
}
