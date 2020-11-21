package httpproxy

import (
	"context"
	"io"
	"net"
	"net/http"
	"sync"
)

func (h *handler) handleHTTPS(responseWriter http.ResponseWriter, request *http.Request) {
	dialer := net.Dialer{}
	destinationConn, err := dialer.DialContext(h.ctx, "tcp", request.Host)
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

	h.wg.Add(1)
	ctx, cancel := context.WithCancel(h.ctx)
	const transferGoroutines = 2
	wg := &sync.WaitGroup{}
	wg.Add(transferGoroutines)
	go func() { // trigger cleanup when done
		wg.Wait()
		cancel()
	}()
	go func() { // cleanup
		<-ctx.Done()
		destinationConn.Close()
		clientConnection.Close()
		h.wg.Done()
	}()
	go transfer(destinationConn, clientConnection, wg)
	go transfer(clientConnection, destinationConn, wg)
}

func transfer(destination io.WriteCloser, source io.ReadCloser, wg *sync.WaitGroup) {
	_, _ = io.Copy(destination, source)
	_ = source.Close()
	_ = destination.Close()
	wg.Done()
}
