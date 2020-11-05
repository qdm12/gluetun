package server

import (
	"fmt"
	"net/http"

	"github.com/qdm12/gluetun/internal/dns"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/openvpn"
	"github.com/qdm12/gluetun/internal/updater"
	"github.com/qdm12/golibs/logging"
)

func newHandler(logger logging.Logger, logging bool,
	buildInfo models.BuildInformation,
	openvpnLooper openvpn.Looper,
	unboundLooper dns.Looper,
	updaterLooper updater.Looper,
) http.Handler {
	return &handler{
		logger:        logger,
		logging:       logging,
		buildInfo:     buildInfo,
		openvpnLooper: openvpnLooper,
		unboundLooper: unboundLooper,
		updaterLooper: updaterLooper,
	}
}

type handler struct {
	logger        logging.Logger
	logging       bool
	buildInfo     models.BuildInformation
	openvpnLooper openvpn.Looper
	unboundLooper dns.Looper
	updaterLooper updater.Looper
}

func (h *handler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	if h.logging {
		h.logger.Info("HTTP %s %s", request.Method, request.RequestURI)
	}
	switch request.Method {
	case http.MethodGet:
		switch request.RequestURI {
		case "/version":
			h.getVersion(responseWriter)
			responseWriter.WriteHeader(http.StatusOK)
		case "/openvpn/actions/restart":
			h.openvpnLooper.Restart()
			responseWriter.WriteHeader(http.StatusOK)
		case "/unbound/actions/restart":
			h.unboundLooper.Restart()
			responseWriter.WriteHeader(http.StatusOK)
		case "/openvpn/portforwarded":
			h.getPortForwarded(responseWriter)
		case "/openvpn/settings":
			h.getOpenvpnSettings(responseWriter)
		case "/updater/restart":
			h.updaterLooper.Restart()
			responseWriter.WriteHeader(http.StatusOK)
		default:
			errString := fmt.Sprintf("Nothing here for %s %s", request.Method, request.RequestURI)
			http.Error(responseWriter, errString, http.StatusBadRequest)
		}
	default:
		errString := fmt.Sprintf("Nothing here for %s %s", request.Method, request.RequestURI)
		http.Error(responseWriter, errString, http.StatusBadRequest)
	}
}
