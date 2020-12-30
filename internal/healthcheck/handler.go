package healthcheck

import (
	"errors"
	"net/http"
	"sync"

	"github.com/qdm12/golibs/logging"
)

type handler struct {
	logger      logging.Logger
	healthErr   error
	healthErrMu sync.RWMutex
}

var errHealthcheckNotRunYet = errors.New("healthcheck did not run yet")

func newHandler(logger logging.Logger) *handler {
	return &handler{
		logger:    logger,
		healthErr: errHealthcheckNotRunYet,
	}
}

func (h *handler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(responseWriter, "method not supported for healthcheck", http.StatusBadRequest)
		return
	}
	if err := h.getErr(); err != nil {
		h.logger.Error(err)
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
	responseWriter.WriteHeader(http.StatusOK)
}

func (h *handler) setErr(err error) {
	h.healthErrMu.Lock()
	defer h.healthErrMu.Unlock()
	h.healthErr = err
}

func (h *handler) getErr() (err error) {
	h.healthErrMu.RLock()
	defer h.healthErrMu.RUnlock()
	return h.healthErr
}
