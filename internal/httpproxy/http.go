package httpproxy

import (
	"context"
	"io"
	"net/http"
)

func (h *handler) handleHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	switch request.URL.Scheme {
	case "http", "https":
	default:
		h.logger.Warn("Unsupported scheme %q", request.URL.Scheme)
		http.Error(responseWriter, "unsupported scheme", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(h.ctx, h.relayTimeout)
	defer cancel()
	request = request.WithContext(ctx)

	request.RequestURI = ""

	for _, key := range hopHeaders {
		request.Header.Del(key)
	}

	if !h.stealth {
		setForwardedHeaders(request)
	}

	response, err := h.client.Do(request)
	if err != nil {
		http.Error(responseWriter, "server error", http.StatusInternalServerError)
		h.logger.Warn("cannot request %s for client %q: %s",
			request.URL, request.RemoteAddr, err)
		return
	}
	defer response.Body.Close()
	if h.verbose {
		h.logger.Info("%s %s %s %s", request.RemoteAddr, response.Status, request.Method, request.URL)
	}

	for _, key := range hopHeaders {
		response.Header.Del(key)
	}

	targetHeaderPtr := responseWriter.Header()
	for key, values := range response.Header {
		for _, value := range values {
			targetHeaderPtr.Add(key, value)
		}
	}

	responseWriter.WriteHeader(response.StatusCode)
	if _, err := io.Copy(responseWriter, response.Body); err != nil {
		h.logger.Error("%s %s: body copy error: %s", request.RemoteAddr, request.URL, err)
	}
}
