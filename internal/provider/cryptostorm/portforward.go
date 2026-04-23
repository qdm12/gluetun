package cryptostorm

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/qdm12/gluetun/internal/provider/common"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

// portRangePattern matches the valid Cryptostorm port range 30000-65535.
const portRangePattern = `([3-5]\d{4}|6[0-4]\d{3}|65[0-4]\d{2}|655[0-2]\d|6553[0-5])`

// regexForwardPlainText matches plain text responses (e.g. from curl):
//
//	37.120.234.253:55555 -> 10.10.123.139:55555
var regexForwardPlainText = regexp.MustCompile(
	`\d+\.\d+\.\d+\.\d+:` + portRangePattern + `\s*->\s*\d+\.\d+\.\d+\.\d+:\d+`)

// regexForwardHTML matches the HTML response from the port forwarding page.
// Each forwarded port has a hidden delete input:
//
//	<input type="hidden" name="delfwd" value="30000">
var regexForwardHTML = regexp.MustCompile(
	`name="delfwd"\s+value="` + portRangePattern + `"`)

// portForwardData is the data persisted to the port forward JSON file.
type portForwardData struct {
	Ports []uint16 `json:"ports"`
}

// PortForward registers a forwarded port with the Cryptostorm port forwarding server
// and returns the active forwarded ports. The server returns plain text listing
// current forwardings. We POST the desired port and parse the response.
// Valid port range is 30000-65535.
// See: https://cryptostorm.is/portfwd
func (p *Provider) PortForward(ctx context.Context, objects utils.PortForwardObjects) (
	internalToExternalPorts map[uint16]uint16, err error,
) {
	const timeout = 10 * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Determine the port to request:
	// 1. Use VPN_PORT_FORWARDING_LISTENING_PORTS[0] if set (non-zero).
	// 2. Otherwise try to read a previously persisted port.
	// 3. Otherwise return an error (Cryptostorm does not auto-assign ports).
	var listeningPort uint16
	if len(objects.ListeningPorts) > 0 && objects.ListeningPorts[0] != 0 {
		listeningPort = objects.ListeningPorts[0]
	}
	if listeningPort == 0 {
		data, err := readPortForwardData(p.portForwardPath)
		if err != nil {
			return nil, fmt.Errorf("reading persisted port forward data: %w", err)
		}
		if len(data.Ports) > 0 {
			listeningPort = data.Ports[0]
		}
	}

	if listeningPort == 0 {
		return nil, fmt.Errorf("%w: set VPN_PORT_FORWARDING_LISTENING_PORTS to a value between 30000 and 65535",
			common.ErrPortForwardNotSupported)
	}

	postBody := "port=" + strconv.FormatUint(uint64(listeningPort), 10)

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

	// The server response lists all currently active forwards for this session,
	// which may include stale ports from prior runs. Only return the port we
	// actually requested so the caller's ListeningPorts slice stays in sync.
	const base, bitSize = 10, 16
	requestedFound := false
	for _, match := range matches {
		portUint64, err := strconv.ParseUint(match[1], base, bitSize)
		if err != nil {
			return nil, fmt.Errorf("parsing port number %q: %w", match[1], err)
		}
		if uint16(portUint64) == listeningPort {
			requestedFound = true
			break
		}
	}
	if !requestedFound {
		return nil, fmt.Errorf("%w: requested port %d not found in server response",
			common.ErrPortForwardNotSupported, listeningPort)
	}

	internalToExternalPorts = map[uint16]uint16{listeningPort: listeningPort}

	// Persist so the next restart can reuse this port without re-requesting.
	if err := writePortForwardData(p.portForwardPath, portForwardData{Ports: []uint16{listeningPort}}); err != nil {
		return nil, fmt.Errorf("persisting port forward data: %w", err)
	}

	return internalToExternalPorts, nil
}

func (p *Provider) KeepPortForward(ctx context.Context,
	_ utils.PortForwardObjects,
) (err error) {
	// Cryptostorm port assignments persist for the session; no keepalive needed.
	<-ctx.Done()
	return ctx.Err()
}

func readPortForwardData(path string) (data portForwardData, err error) {
	file, err := os.Open(path)
	if os.IsNotExist(err) {
		return data, nil
	} else if err != nil {
		return data, err
	}

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		_ = file.Close()
		return data, err
	}

	return data, file.Close()
}

func writePortForwardData(path string, data portForwardData) (err error) {
	const dirPermission = fs.FileMode(0o755)
	if err := os.MkdirAll(filepath.Dir(path), dirPermission); err != nil {
		return err
	}

	const permission = fs.FileMode(0o644)
	file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, permission)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(data); err != nil {
		_ = file.Close()
		return err
	}

	return file.Close()
}
