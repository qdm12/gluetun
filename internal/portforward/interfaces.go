package portforward

import (
	"context"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/models"
)

type PortAllower interface {
	SetAllowedPort(ctx context.Context, port uint16, intf string) (err error)
	RemoveAllowedPort(ctx context.Context, port uint16) (err error)
}

type statusManager interface {
	GetStatus() (status models.LoopStatus)
	SetStatus(status models.LoopStatus)
	ApplyStatus(ctx context.Context, status models.LoopStatus) (
		outcome string, err error)
}

type StateManager interface {
	GetSettings() (settings settings.PortForwarding)
	SetSettings(ctx context.Context,
		settings settings.PortForwarding) (outcome string)
	GetPortForwarded() (port uint16)
	SetPortForwarded(port uint16)
	GetStartData() (startData StartData)
	SetStartData(startData StartData)
}
