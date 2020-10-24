package httpproxy

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/qdm12/golibs/logging"
)

func newHandler(ctx context.Context, wg *sync.WaitGroup,
	client *http.Client, logger logging.Logger,
	stealth, verbose bool) http.Handler {
	const relayTimeout = 10 * time.Second
	return &handler{
		ctx:          ctx,
		wg:           wg,
		client:       client,
		logger:       logger,
		relayTimeout: relayTimeout,
		verbose:      verbose,
		stealth:      stealth,
	}
}

type handler struct {
	ctx          context.Context
	wg           *sync.WaitGroup
	client       *http.Client
	logger       logging.Logger
	relayTimeout time.Duration
	verbose      bool
	stealth      bool
}

func (h *handler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodConnect:
		h.handleHTTPS(responseWriter, request)
	default:
		h.handleHTTP(responseWriter, request)
	}
}

func setForwardedHeaders(request *http.Request) {
	clientIP, _, err := net.SplitHostPort(request.RemoteAddr)
	if err != nil {
		return
	}
	// keep existing proxy headers
	if prior, ok := request.Header["X-Forwarded-For"]; ok {
		clientIP = fmt.Sprintf("%s,%s", strings.Join(prior, ", "), clientIP)
	}
	request.Header.Set("X-Forwarded-For", clientIP)
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
