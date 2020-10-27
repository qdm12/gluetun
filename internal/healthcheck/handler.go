package healthcheck

import (
	"net"
	"net/http"

	"github.com/qdm12/golibs/logging"
)

type handler struct {
	logger   logging.Logger
	resolver *net.Resolver
}

func newHandler(logger logging.Logger, resolver *net.Resolver) http.Handler {
	return &handler{
		logger:   logger,
		resolver: resolver,
	}
}

func (h *handler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(responseWriter, "method not supported for healthcheck", http.StatusBadRequest)
		return
	}
	err := healthCheck(request.Context(), h.resolver)
	if err != nil {
		h.logger.Error(err)
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
	responseWriter.WriteHeader(http.StatusOK)
}
