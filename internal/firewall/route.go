package firewall

import (
	"encoding/hex"
	"net"

	"fmt"
	"strings"

	"github.com/qdm12/private-internet-access-docker/internal/constants"
)

func (c *configurator) AddRoutesVia(subnets []net.IPNet, defaultGateway net.IP, defaultInterface string) error {
	for _, subnet := range subnets {
		c.logger.Info("adding %s as route via %s", subnet, defaultInterface)
		_, err := c.commander.Run("ip", "route", "add", subnet.String(), "via", defaultGateway.String(), "dev", defaultInterface)
		if err != nil {
			return fmt.Errorf("cannot add route for %s: %w", subnet, err)
		}
	}
	return nil
}

func (c *configurator) GetDefaultRoute() (defaultInterface string, defaultGateway net.IP, defaultSubnet net.IPNet, err error) {
	data, err := c.fileManager.ReadFile(string(constants.NetRoute))
	if err != nil {
		return "", nil, defaultSubnet, err
	}
	// Verify number of lines and fields
	lines := strings.Split(string(data), "\n")
	if len(lines) < 3 {
		return "", nil, defaultSubnet, fmt.Errorf("not enough lines (%d) found in %s", len(lines), constants.NetRoute)
	}
	fieldsLine1 := strings.Fields(lines[1])
	if len(fieldsLine1) < 3 {
		return "", nil, defaultSubnet, fmt.Errorf("not enough fields in %q", lines[1])
	}
	fieldsLine2 := strings.Fields(lines[2])
	if len(fieldsLine2) < 8 {
		return "", nil, defaultSubnet, fmt.Errorf("not enough fields in %q", lines[2])
	}
	// get information
	defaultInterface = fieldsLine1[0]
	defaultGateway, err = reversedHexToIPv4(fieldsLine1[2])
	if err != nil {
		return "", nil, defaultSubnet, err
	}
	netNumber, err := reversedHexToIPv4(fieldsLine2[1])
	if err != nil {
		return "", nil, defaultSubnet, err
	}
	netMask, err := hexToIPv4Mask(fieldsLine2[7])
	if err != nil {
		return "", nil, defaultSubnet, err
	}
	return defaultInterface, defaultGateway, net.IPNet{IP: netNumber, Mask: netMask}, nil
}

func reversedHexToIPv4(reversedHex string) (IP net.IP, err error) {
	bytes, err := hex.DecodeString(reversedHex)
	if err != nil {
		return nil, fmt.Errorf("cannot parse reversed IP hex %q: %s", reversedHex, err)
	} else if len(bytes) != 4 {
		return nil, fmt.Errorf("hex string contains %d bytes instead of 4", len(bytes))
	}
	return []byte{bytes[3], bytes[2], bytes[1], bytes[0]}, nil
}

func hexToIPv4Mask(hexString string) (mask net.IPMask, err error) {
	bytes, err := hex.DecodeString(hexString)
	if err != nil {
		return nil, fmt.Errorf("cannot parse hex mask %q: %s", hexString, err)
	} else if len(bytes) != 4 {
		return nil, fmt.Errorf("hex string contains %d bytes instead of 4", len(bytes))
	}
	return []byte{bytes[3], bytes[2], bytes[1], bytes[0]}, nil
}
