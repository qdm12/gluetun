package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/golibs/logging"
)

func newHandlerV1(logger logging.Logger, buildInfo models.BuildInformation,
	openvpn, dns, updater http.Handler) http.Handler {
	return &handlerV1{
		logger:    logger,
		buildInfo: buildInfo,
		openvpn:   openvpn,
		dns:       dns,
		updater:   updater,
	}
}

type handlerV1 struct {
	logger    logging.Logger
	buildInfo models.BuildInformation
	openvpn   http.Handler
	dns       http.Handler
	updater   http.Handler
}

func (h *handlerV1) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.RequestURI == "/version" && r.Method == http.MethodGet:
		h.getVersion(w)
	case strings.HasPrefix(r.RequestURI, "/openvpn"):
		h.openvpn.ServeHTTP(w, r)
	case strings.HasPrefix(r.RequestURI, "/dns"):
		h.dns.ServeHTTP(w, r)
	case strings.HasPrefix(r.RequestURI, "/updater"):
		h.updater.ServeHTTP(w, r)
	default:
		errString := fmt.Sprintf("%s %s not found", r.Method, r.RequestURI)
		http.Error(w, errString, http.StatusNotFound)
	}
}

func (h *handlerV1) getVersion(w http.ResponseWriter) {
	data, err := json.Marshal(h.buildInfo)
	if err != nil {
		h.logger.Warn(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if _, err := w.Write(data); err != nil {
		h.logger.Warn(err)
	}
}
