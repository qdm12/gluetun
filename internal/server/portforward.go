package server

import (
	"context"
	"encoding/json"
	"net/http"
)

func newPortForwardHandler(ctx context.Context,
	portForward PortForwardedGetter, warner warner,
) http.Handler {
	return &portForwardHandler{
		ctx:         ctx,
		portForward: portForward,
		warner:      warner,
	}
}

type portForwardHandler struct {
	ctx         context.Context //nolint:containedctx
	portForward PortForwardedGetter
	warner      warner
}

func (h *portForwardHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getPortForwarded(w)
	default:
		errMethodNotSupported(w, r.Method)
	}
}

func (h *portForwardHandler) getPortForwarded(w http.ResponseWriter) {
	ports := h.portForward.GetPortsForwarded()
	encoder := json.NewEncoder(w)
	var data any
	switch len(ports) {
	case 0:
		data = portWrapper{Port: 0} // TODO v4 change to portsWrapper
	case 1:
		data = portWrapper{Port: ports[0]} // TODO v4 change to portsWrapper
	default:
		data = portsWrapper{Ports: ports}
	}

	err := encoder.Encode(data)
	if err != nil {
		h.warner.Warn(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
}
