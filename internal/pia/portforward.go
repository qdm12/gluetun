package pia

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/qdm12/golibs/files"
	"github.com/qdm12/private-internet-access-docker/internal/constants"
	"github.com/qdm12/private-internet-access-docker/internal/models"
)

func (c *configurator) GetPortForward() (port uint16, err error) {
	b, err := c.random.GenerateRandomBytes(32)
	if err != nil {
		return 0, err
	}
	clientID := hex.EncodeToString(b)
	url := fmt.Sprintf("%s/?client_id=%s", constants.PIAPortForwardURL, clientID)
	content, status, err := c.client.GetContent(url)
	switch {
	case err != nil:
		return 0, err
	case status != http.StatusOK:
		return 0, fmt.Errorf("status is %d for %s; does your PIA server support port forwarding?", status, url)
	case len(content) == 0:
		return 0, fmt.Errorf("port forwarding is already activated on this connection, has expired, or you are not connected to a PIA region that supports port forwarding")
	}
	body := struct {
		Port uint16 `json:"port"`
	}{}
	if err := json.Unmarshal(content, &body); err != nil {
		return 0, fmt.Errorf("port forwarding response: %w", err)
	}
	return body.Port, nil
}

func (c *configurator) WritePortForward(filepath models.Filepath, port uint16, uid, gid int) (err error) {
	return c.fileManager.WriteLinesToFile(
		string(filepath),
		[]string{fmt.Sprintf("%d", port)},
		files.Ownership(uid, gid),
		files.Permissions(0400))
}

func (c *configurator) AllowPortForwardFirewall(ctx context.Context, device models.VPNDevice, port uint16) (err error) {
	return c.firewall.AllowInputTrafficOnPort(ctx, device, port)
}
