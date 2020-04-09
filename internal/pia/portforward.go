package pia

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/qdm12/golibs/files"
	"github.com/qdm12/private-internet-access-docker/internal/constants"
	"github.com/qdm12/private-internet-access-docker/internal/models"
)

func (c *configurator) GetPortForward() (port uint16, err error) {
	c.logger.Info("%s: Obtaining port to be forwarded", logPrefix)
	b, err := c.random.GenerateRandomBytes(32)
	if err != nil {
		return 0, err
	}
	clientID := hex.EncodeToString(b)
	url := fmt.Sprintf("%s/?client_id=%s", constants.PIAPortForwardURL, clientID)
	content, status, err := c.client.GetContent(url)
	if err != nil {
		return 0, err
	} else if status != 200 {
		return 0, fmt.Errorf("status is %d for %s; does your PIA server support port forwarding?", status, url)
	} else if len(content) == 0 {
		return 0, fmt.Errorf("port forwarding is already activated on this connection, has expired, or you are not connected to a PIA region that supports port forwarding")
	}
	body := struct {
		Port uint16 `json:"port"`
	}{}
	if err := json.Unmarshal(content, &body); err != nil {
		return 0, fmt.Errorf("port forwarding response: %w", err)
	}
	c.logger.Info("%s: Port forwarded is %d", logPrefix, body.Port)
	return body.Port, nil
}

func (c *configurator) WritePortForward(filepath models.Filepath, port uint16, uid, gid int) (err error) {
	c.logger.Info("%s: Writing forwarded port to %s", logPrefix, filepath)
	return c.fileManager.WriteLinesToFile(
		string(filepath),
		[]string{fmt.Sprintf("%d", port)},
		files.Ownership(uid, gid),
		files.Permissions(400))
}

func (c *configurator) AllowPortForwardFirewall(device models.VPNDevice, port uint16) (err error) {
	c.logger.Info("%s: Allowing forwarded port %d through firewall", logPrefix, port)
	return c.firewall.AllowInputTrafficOnPort(device, port)
}

func (c *configurator) ClearPortForward(filepath models.Filepath, uid, gid int) (err error) {
	c.logger.Info("%s: Clearing forwarded port status file %s", logPrefix, filepath)
	return c.fileManager.WriteToFile(string(filepath), nil, files.Ownership(uid, gid), files.Permissions(400))
}
