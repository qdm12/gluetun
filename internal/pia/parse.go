package pia

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/qdm12/private-internet-access-docker/internal/constants"
	"github.com/qdm12/private-internet-access-docker/internal/models"
)

func (c *configurator) ParseConfig(lines []string) (IPs []net.IP, port uint16, device models.VPNDevice, err error) {
	c.logger.Info("%s: parsing openvpn configuration", logPrefix)
	remoteLineFound := false
	deviceLineFound := false
	for _, line := range lines {
		if strings.HasPrefix(line, "remote ") {
			remoteLineFound = true
			words := strings.Fields(line)
			if len(words) != 3 {
				return nil, 0, "", fmt.Errorf("line %q misses information", line)
			}
			host := words[1]
			if err := c.verifyPort(words[2]); err != nil {
				return nil, 0, "", fmt.Errorf("line %q has an invalid port: %w", line, err)
			}
			portUint64, err := strconv.ParseUint(words[2], 10, 16)
			if err != nil {
				return nil, 0, "", err
			}
			port = uint16(portUint64)
			IPs, err = c.lookupIP(host)
			if err != nil {
				return nil, 0, "", err
			}
		} else if strings.HasPrefix(line, "dev ") {
			deviceLineFound = true
			fields := strings.Fields(line)
			if len(fields) != 2 {
				return nil, 0, "", fmt.Errorf("line %q misses information", line)
			}
			device = models.VPNDevice(fields[1] + "0")
			if device != constants.TUN && device != constants.TAP {
				return nil, 0, "", fmt.Errorf("device %q is not valid", device)
			}
		}
	}
	if remoteLineFound && deviceLineFound {
		c.logger.Info("%s: Found %d PIA server IP addresses, port %d and device %s", logPrefix, len(IPs), port, device)
		return IPs, port, device, nil
	} else if !remoteLineFound {
		return nil, 0, "", fmt.Errorf("remote line not found in Openvpn configuration")
	} else {
		return nil, 0, "", fmt.Errorf("device line not found in Openvpn configuration")
	}
}
