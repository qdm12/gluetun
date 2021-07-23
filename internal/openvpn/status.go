package openvpn

import (
	"context"

	"github.com/qdm12/gluetun/internal/models"
)

func (l *looper) GetStatus() (status models.LoopStatus) {
	return l.statusManager.GetStatus()
}

func (l *looper) ApplyStatus(ctx context.Context, status models.LoopStatus) (
	outcome string, err error) {
	return l.statusManager.ApplyStatus(ctx, status)
}
