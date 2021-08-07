package dns

import (
	"context"

	"github.com/qdm12/gluetun/internal/models"
)

func (l *Loop) GetStatus() (status models.LoopStatus) {
	return l.statusManager.GetStatus()
}

func (l *Loop) ApplyStatus(ctx context.Context, status models.LoopStatus) (
	outcome string, err error) {
	return l.statusManager.ApplyStatus(ctx, status)
}
