package privatevpn

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"

	"github.com/qdm12/gluetun/internal/provider/common"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

var regexPort = regexp.MustCompile(`[1-9][0-9]{0,4}`)

var ErrPortForwardedNotFound = errors.New("port forwarded not found")

// PortForward obtains a VPN server side port forwarded from the PrivateVPN API.
// It returns 0 if all ports are to forwarded on a dedicated server IP.
func (p *Provider) PortForward(ctx context.Context, objects utils.PortForwardObjects) (
	ports []uint16, err error,
) {
	url := "https://connect.pvdatanet.com/v3/Api/port?ip[]=" + objects.InternalIP.String()
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating HTTP request: %w", err)
	}

	response, err := objects.Client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("sending HTTP request: %w", err)
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %d %s", common.ErrHTTPStatusCodeNotOK,
			response.StatusCode, response.Status)
	}

	defer response.Body.Close()

	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	var data struct {
		Status    string `json:"status"`
		Supported bool   `json:"supported"`
	}
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, fmt.Errorf("decoding JSON response: %w; data is: %s",
			err, string(bytes))
	} else if !data.Supported {
		return nil, fmt.Errorf("%w for this VPN server", common.ErrPortForwardNotSupported)
	}

	portString := regexPort.FindString(data.Status)
	if portString == "" {
		return nil, fmt.Errorf("%w: in status %q", ErrPortForwardedNotFound, data.Status)
	}

	const base, bitSize = 10, 16
	portUint64, err := strconv.ParseUint(portString, base, bitSize)
	if err != nil {
		return nil, fmt.Errorf("parsing port: %w", err)
	}
	return []uint16{uint16(portUint64)}, nil
}

func (p *Provider) KeepPortForward(ctx context.Context,
	_ utils.PortForwardObjects,
) (err error) {
	<-ctx.Done()
	return ctx.Err()
}
