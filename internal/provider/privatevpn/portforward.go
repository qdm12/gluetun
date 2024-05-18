package privatevpn

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/qdm12/gluetun/internal/provider/common"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

var (
	regexPort = regexp.MustCompile(`^[1-9][0-9]{0,4}$`)
)

var (
	ErrPortForwardedNotFound = errors.New("port forwarded not found")
)

// PortForward obtains a VPN server side port forwarded from the PrivateVPN API.
// It returns 0 if all ports are to forwarded on a dedicated server IP.
func (p *Provider) PortForward(ctx context.Context, objects utils.PortForwardObjects) (
	port uint16, err error) {
	url := "https://connect.pvdatanet.com/v3/Api/port?ip[]=" + objects.ServerIP.String()
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, fmt.Errorf("creating HTTP request: %w", err)
	}

	response, err := objects.Client.Do(request)
	if err != nil {
		return 0, fmt.Errorf("sending HTTP request: %w", err)
	}

	if response.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("%w: %d %s", common.ErrHTTPStatusCodeNotOK,
			response.StatusCode, response.Status)
	}

	defer response.Body.Close()
	decoder := json.NewDecoder(response.Body)
	var data struct {
		Status    string `json:"status"`
		Supported bool   `json:"supported"`
	}
	err = decoder.Decode(&data)
	if err != nil {
		return 0, fmt.Errorf("decoding JSON response: %w", err)
	} else if !data.Supported {
		return 0, fmt.Errorf("%w: for server IP %s",
			common.ErrPortForwardNotSupported, objects.ServerIP)
	}

	portString := regexPort.FindString(data.Status)
	if portString == "" {
		return 0, fmt.Errorf("%w: in status %q", ErrPortForwardedNotFound, data.Status)
	}

	const base, bitSize = 10, 16
	portUint64, err := strconv.ParseUint(portString, base, bitSize)
	if err != nil {
		return 0, fmt.Errorf("parsing port %q: %w", portString, err)
	}
	port = uint16(portUint64)
	return port, nil
}

func (p *Provider) KeepPortForward(ctx context.Context,
	_ utils.PortForwardObjects) (err error) {
	<-ctx.Done()
	return ctx.Err()
}
