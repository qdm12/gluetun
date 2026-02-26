package cryptostorm

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/qdm12/gluetun/internal/provider/common"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

// regexActivePort matches forwarded ports listed in the Cryptostorm HTML response,
// e.g. <p class="list-group-item-header">Port 55555</p>
// Valid port range per Cryptostorm is 30000-65535.
var regexActivePort = regexp.MustCompile(`list-group-item-header">Port ([0-9]+)<`)

// cryptostormPFServer is the fixed internal IP of Cryptostorm's
// port forwarding server, reachable only from within the VPN tunnel.
const cryptostormPFServer = "10.31.33.7"

// PortForward registers a forwarded port with the Cryptostorm port forwarding server
// and returns the active forwarded ports. The server always returns an HTML page;
// we POST the desired port and then parse the current forwards list from the response.
// If the port is already forwarded (e.g. from a previous session) it will appear in
// the list regardless of whether the POST succeeded, so we treat that as success.
// Valid port range is 30000–65535.
// See: https://cryptostorm.is/port_forwarding
func (p *Provider) PortForward(ctx context.Context, objects utils.PortForwardObjects) (
	ports []uint16, err error,
) {
	const timeout = 10 * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	postBody := ""
	if objects.ListeningPort != 0 {
		postBody = "port=" + strconv.FormatUint(uint64(objects.ListeningPort), 10)
	}

	pfURL := "http://" + cryptostormPFServer + "/fwd"
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, pfURL,
		strings.NewReader(postBody))
	if err != nil {
		return nil, fmt.Errorf("creating HTTP request: %w", err)
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response, err := objects.Client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("sending HTTP request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %d %s", common.ErrHTTPStatusCodeNotOK,
			response.StatusCode, response.Status)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	// Parse all currently active port forwards from the HTML response.
	matches := regexActivePort.FindAllStringSubmatch(string(body), -1)
	if len(matches) == 0 {
		return nil, fmt.Errorf("%w: no active port forwards found in response",
			common.ErrPortForwardNotSupported)
	}

	const base, bitSize = 10, 16
	for _, match := range matches {
		portUint64, err := strconv.ParseUint(match[1], base, bitSize)
		if err != nil {
			return nil, fmt.Errorf("parsing port number %q: %w", match[1], err)
		}
		ports = append(ports, uint16(portUint64))
	}

	return ports, nil
}

func (p *Provider) KeepPortForward(ctx context.Context,
	_ utils.PortForwardObjects,
) (err error) {
	// Cryptostorm port assignments persist for the session; no keepalive needed.
	<-ctx.Done()
	return ctx.Err()
}
