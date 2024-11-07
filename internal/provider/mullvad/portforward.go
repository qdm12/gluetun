package mullvad

import (
	"context"

	"github.com/qdm12/gluetun/internal/provider/utils"
)

// PortForward obtains a VPN server side port forwarded from ProtonVPN gateway.
func (p *Provider) PortForward(_ context.Context, objects utils.PortForwardObjects) (
	port uint16, err error) {
	objects.Logger.Debug("mullvad: port forward")
	port = 10000
	return port, nil
}

func (p *Provider) KeepPortForward(ctx context.Context,
	objects utils.PortForwardObjects) (err error) {
	objects.Logger.Debug("mullvad: keeping port forward")
	<-ctx.Done()
	objects.Logger.Debug("mullvad: keeping port forward exiting")
	return nil
}
