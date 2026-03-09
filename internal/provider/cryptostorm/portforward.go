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

// regexForwardPlainText matches plain text responses (e.g. from curl):
//
//	37.120.234.253:55555 -> 10.10.123.139:55555
var regexForwardPlainText = regexp.MustCompile(
	`\d+\.\d+\.\d+\.\d+:(\d+)\s*->\s*\d+\.\d+\.\d+\.\d+:\d+`)

// regexForwardHTML matches the HTML response from the port forwarding page.
// Each forwarded port has a hidden delete input:
//
//	<input type="hidden" name="delfwd" value="30000">
var regexForwardHTML = regexp.MustCompile(
	`name="delfwd"\s+value="(\d+)"`)

// PortForward registers a forwarded port with the Cryptostorm port forwarding server
// and returns the active forwarded ports. The server returns plain text listing
// current forwardings. We POST the desired port and parse the response.
// Valid port range is 30000-65535.
// See: https://cryptostorm.is/portfwd
func (p *Provider) PortForward(ctx context.Context, objects utils.PortForwardObjects) (
	ports []uint16, err error,
) {
	const timeout = 10 * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Cryptostorm requires a port in the 30000-65535 range to be specified
	// via VPN_PORT_FORWARDING_LISTENING_PORT.
	if objects.ListeningPort == 0 {
		return nil, fmt.Errorf("%w: set VPN_PORT_FORWARDING_LISTENING_PORT to a value between 30000 and 65535",
			common.ErrPortForwardNotSupported)
	}
	postBody := "port=" + strconv.FormatUint(uint64(objects.ListeningPort), 10)

	// IPv4: http://10.31.33.7/fwd
	// IPv6: http://[2001:db8::7]/fwd (for future use)
	const portForwardURL = "http://10.31.33.7/fwd"
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, portForwardURL,
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

	// Parse forwarded ports from the response. The server returns HTML to
	// Go's HTTP client but plain text to curl, so we try both formats.
	bodyStr := string(body)
	matches := regexForwardPlainText.FindAllStringSubmatch(bodyStr, -1)
	if len(matches) == 0 {
		matches = regexForwardHTML.FindAllStringSubmatch(bodyStr, -1)
	}
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
