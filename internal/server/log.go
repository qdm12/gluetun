package server

import (
	"net/http"
	"sync"
	"time"

	"github.com/qdm12/golibs/logging"
)

func withLogMiddleware(childHandler http.Handler, logger logging.Logger, enabled bool) *logMiddleware {
	return &logMiddleware{
		childHandler: childHandler,
		logger:       logger,
		timeNow:      time.Now,
		enabled:      enabled,
	}
}

type logMiddleware struct {
	childHandler http.Handler
	logger       logging.Logger
	timeNow      func() time.Time
	enabled      bool
	enabledMu    sync.RWMutex
}

func (m *logMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !m.isEnabled() {
		m.childHandler.ServeHTTP(w, r)
		return
	}
	tStart := m.timeNow()
	statefulWriter := &statefulResponseWriter{httpWriter: w}
	m.childHandler.ServeHTTP(statefulWriter, r)
	duration := m.timeNow().Sub(tStart)
	m.logger.Info("%d %s %s wrote %dB to %s in %s",
		statefulWriter.statusCode, r.Method, r.RequestURI, statefulWriter.length, r.RemoteAddr, duration)
}

func (m *logMiddleware) setEnabled(enabled bool) {
	m.enabledMu.Lock()
	defer m.enabledMu.Unlock()
	m.enabled = enabled
}

func (m *logMiddleware) isEnabled() (enabled bool) {
	m.enabledMu.RLock()
	defer m.enabledMu.RUnlock()
	return m.enabled
}

type statefulResponseWriter struct {
	httpWriter http.ResponseWriter
	statusCode int
	length     int
}

func (w *statefulResponseWriter) Write(b []byte) (n int, err error) {
	n, err = w.httpWriter.Write(b)
	if w.statusCode == 0 {
		w.statusCode = http.StatusOK
	}
	w.length += n
	return n, err
}

func (w *statefulResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.httpWriter.WriteHeader(statusCode)
}

func (w *statefulResponseWriter) Header() http.Header {
	return w.httpWriter.Header()
}
