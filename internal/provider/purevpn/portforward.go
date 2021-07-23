package purevpn

import (
	"context"
	"net"
	"net/http"

	"github.com/qdm12/gluetun/internal/firewall"
	"github.com/qdm12/golibs/logging"
)

func (p *Purevpn) PortForward(ctx context.Context, client *http.Client,
	pfLogger logging.Logger, gateway net.IP, fw firewall.Configurator,
	syncState func(port uint16) (pfFilepath string)) {
	panic("port forwarding is not supported for PureVPN")
}
