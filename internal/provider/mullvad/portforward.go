package mullvad

import (
	"context"
	"net"
	"net/http"

	"github.com/qdm12/gluetun/internal/firewall"
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/os"
)

func (m *Mullvad) PortForward(ctx context.Context, client *http.Client,
	openFile os.OpenFileFunc, pfLogger logging.Logger, gateway net.IP,
	fw firewall.Configurator, syncState func(port uint16) (pfFilepath string)) {
	panic("port forwarding logic is not needed for Mullvad")
}
