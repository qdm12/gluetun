package server

import (
	"context"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/publicip/ipinfo"
)

type VPNLooper interface {
	GetStatus() (status models.LoopStatus)
	ApplyStatus(ctx context.Context, status models.LoopStatus) (
		outcome string, err error)
	GetSettings() (settings settings.VPN)
}

type DNSLoop interface {
	ApplyStatus(ctx context.Context, status models.LoopStatus) (
		outcome string, err error)
	GetStatus() (status models.LoopStatus)
}

type PortForwardedGetter interface {
	GetPortForwarded() (portForwarded uint16)
}

type PublicIPLoop interface {
	GetData() (data ipinfo.Response)
}
