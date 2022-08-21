package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/qdm12/gluetun/internal/models"
)

func newHandlerV1(w warner, buildInfo models.BuildInformation,
	vpn, openvpn, dns, updater, publicip http.Handler) http.Handler {
	return &handlerV1{
		warner:    w,
		buildInfo: buildInfo,
		vpn:       vpn,
		openvpn:   openvpn,
		dns:       dns,
		updater:   updater,
		publicip:  publicip,
	}
}

type handlerV1 struct {
	warner    warner
	buildInfo models.BuildInformation
	vpn       http.Handler
	openvpn   http.Handler
	dns       http.Handler
	updater   http.Handler
	publicip  http.Handler
}

func (h *handlerV1) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.RequestURI == "/version" && r.Method == http.MethodGet:
		h.getVersion(w)
	case strings.HasPrefix(r.RequestURI, "/vpn"):
		h.vpn.ServeHTTP(w, r)
	case strings.HasPrefix(r.RequestURI, "/openvpn"):
		h.openvpn.ServeHTTP(w, r)
	case strings.HasPrefix(r.RequestURI, "/dns"):
		h.dns.ServeHTTP(w, r)
	case strings.HasPrefix(r.RequestURI, "/updater"):
		h.updater.ServeHTTP(w, r)
	case strings.HasPrefix(r.RequestURI, "/publicip"):
		h.publicip.ServeHTTP(w, r)
	default:
		errString := fmt.Sprintf("%s %s not found", r.Method, r.RequestURI)
		http.Error(w, errString, http.StatusBadRequest)
	}
}

func (h *handlerV1) getVersion(w http.ResponseWriter) {
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(h.buildInfo); err != nil {
		h.warner.Warn(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
}
