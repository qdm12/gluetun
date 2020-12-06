//nolint:dupl
package server

import (
	"net/http"
	"strings"

	"github.com/qdm12/gluetun/internal/updater"
	"github.com/qdm12/golibs/logging"
)

func newUpdaterHandler(
	looper updater.Looper,
	logger logging.Logger) http.Handler {
	return &updaterHandler{
		looper: looper,
		logger: logger,
	}
}

type updaterHandler struct {
	looper updater.Looper
	logger logging.Logger
}

func (h *updaterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.RequestURI = strings.TrimPrefix(r.RequestURI, "/updater")
	switch r.RequestURI {
	case "/status":
		switch r.Method {
		case http.MethodGet:
			h.getStatus(w)
		case http.MethodPut:
			h.setStatus(w, r)
		default:
			http.Error(w, "", http.StatusNotFound)
		}
	default:
		http.Error(w, "", http.StatusNotFound)
	}
}

func (h *updaterHandler) getStatus(w http.ResponseWriter) {
	// status := h.looper.GetStatus()
	// encoder := json.NewEncoder(w)
	// data := statusWrapper{Status: string(status)}
	// if err := encoder.Encode(data); err != nil {
	// 	h.logger.Warn(err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }
}

func (h *updaterHandler) setStatus(w http.ResponseWriter, r *http.Request) {
	// decoder := json.NewDecoder(r.Body)
	// var data statusWrapper
	// if err := decoder.Decode(&data); err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }
	// status := models.LoopStatus(data.Status)
	// switch status {
	// case constants.Stopped, constants.Running:
	// default:
	// 	errString := fmt.Sprintf(
	// 		"invalid openvpn status %q: possible values are: %s, %s",
	// 		status, constants.Stopped, constants.Running)
	// 	http.Error(w, errString, http.StatusBadRequest)
	// 	return
	// }
	// message := h.looper.SetStatus(status)
	// if _, err := w.Write([]byte(message)); err != nil {
	// 	h.logger.Warn(err)
	// }
}
