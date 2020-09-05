package updater

import (
	"context"
	"net"
	"net/http"
)

type (
	httpGetFunc  func(url string) (r *http.Response, err error)
	lookupIPFunc func(ctx context.Context, host string) (ips []net.IP, err error)
)
