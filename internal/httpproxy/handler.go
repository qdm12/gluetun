package httpproxy

import (
	"context"
	"net/http"
	"sync"
	"time"
)

func newHandler(ctx context.Context, wg *sync.WaitGroup, logger Logger,
	stealth, verbose bool, username, password string) http.Handler {
	const httpTimeout = 24 * time.Hour
	return &handler{
		ctx: ctx,
		wg:  wg,
		client: &http.Client{
			Timeout:       httpTimeout,
			CheckRedirect: returnRedirect},
		logger:   logger,
		verbose:  verbose,
		stealth:  stealth,
		username: username,
		password: password,
	}
}

type handler struct {
	ctx                context.Context //nolint:containedctx
	wg                 *sync.WaitGroup
	client             *http.Client
	logger             Logger
	verbose, stealth   bool
	username, password string
}

func (h *handler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	if !h.isAccepted(responseWriter, request) {
		return
	}
	if !h.isAuthorized(responseWriter, request) {
		return
	}
	request.Header.Del("Proxy-Connection")
	request.Header.Del("Proxy-Authenticate")
	request.Header.Del("Proxy-Authorization")
	switch request.Method {
	case http.MethodConnect:
		h.handleHTTPS(responseWriter, request)
	default:
		h.handleHTTP(responseWriter, request)
	}
}

// http://www.w3.org/Protocols/rfc2616/rfc2616-sec13.html
var hopHeaders = [...]string{ //nolint:gochecknoglobals
	"Connection",
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Te", // canonicalized version of "TE"
	"Trailers",
	"Transfer-Encoding",
	"Upgrade",
}

// Do not follow redirect, but directly return the redirect response.
func returnRedirect(req *http.Request, via []*http.Request) error {
	return http.ErrUseLastResponse
}
