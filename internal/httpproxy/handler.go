package httpproxy

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/qdm12/golibs/logging"
)

func newHandler(ctx context.Context, wg *sync.WaitGroup,
	client *http.Client, logger logging.Logger,
	stealth, verbose bool, username, password string) http.Handler {
	const relayTimeout = 10 * time.Second
	return &handler{
		ctx:          ctx,
		wg:           wg,
		client:       client,
		logger:       logger,
		relayTimeout: relayTimeout,
		verbose:      verbose,
		stealth:      stealth,
		username:     username,
		password:     password,
	}
}

type handler struct {
	ctx                context.Context
	wg                 *sync.WaitGroup
	client             *http.Client
	logger             logging.Logger
	relayTimeout       time.Duration
	verbose, stealth   bool
	username, password string
}

func (h *handler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	if len(h.username) > 0 && !isAuthorized(responseWriter, request, h.username, h.password) {
		h.logger.Info("%s unauthorized", request.RemoteAddr)
		return
	}
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
