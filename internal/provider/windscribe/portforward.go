package windscribe

import (
	"context"
	"net"
	"net/http"

	"github.com/qdm12/gluetun/internal/firewall"
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/os"
)

func (w *Windscribe) PortForward(ctx context.Context, clienw *http.Client,
	openFile os.OpenFileFunc, pfLogger logging.Logger, gateway net.IP,
	fw firewall.Configurator, syncState func(port uint16) (pfFilepath string)) {
	panic("port forwarding is not supported for Windscribe")
}
