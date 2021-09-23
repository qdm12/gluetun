package healthcheck

import (
	"errors"
	"net/http"
	"sync"
)

type handler struct {
	healthErr   error
	healthErrMu sync.RWMutex
}

var errHealthcheckNotRunYet = errors.New("healthcheck did not run yet")

func newHandler() *handler {
	return &handler{
		healthErr: errHealthcheckNotRunYet,
	}
}

func (h *handler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(responseWriter, "method not supported for healthcheck", http.StatusBadRequest)
		return
	}
	if err := h.getErr(); err != nil {
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
