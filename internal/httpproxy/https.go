package httpproxy

import (
	"io"
	"net"
	"net/http"
	"sync"
)

func (h *handler) handleHTTPS(responseWriter http.ResponseWriter, request *http.Request) {
	destinationConn, err := net.DialTimeout("tcp", request.Host, h.relayTimeout)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusServiceUnavailable)
		return
	}

	responseWriter.WriteHeader(http.StatusOK)

	hijacker, ok := responseWriter.(http.Hijacker)
	if !ok {
		http.Error(responseWriter, "Hijacking not supported", http.StatusInternalServerError)
		return
	}
	clientConnection, _, err := hijacker.Hijack()
	if err != nil {
		h.logger.Warn(err)
		http.Error(responseWriter, err.Error(), http.StatusServiceUnavailable)
		if err := destinationConn.Close(); err != nil {
			h.logger.Error("closing destination connection: %s", err)
		}
		return
	}

	if h.verbose {
		h.logger.Info("%s <-> %s", request.RemoteAddr, request.Host)
	}

	ctx := h.ctx
	wg := h.wg
	destinationCtxConn := newContextConn(ctx, destinationConn)
	clientCtxConn := newContextConn(ctx, clientConnection)
	const transferGoroutines = 2
	h.wg.Add(transferGoroutines)
	go transfer(destinationCtxConn, clientCtxConn, wg)
	go transfer(clientCtxConn, destinationCtxConn, wg)
}

func transfer(destination io.WriteCloser, source io.ReadCloser, wg *sync.WaitGroup) {
	_, _ = io.Copy(destination, source)
	_ = source.Close()
	_ = destination.Close()
	wg.Done()
}
