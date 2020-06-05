package cli

import (
	"fmt"
	"net"
	"time"

	"github.com/qdm12/golibs/files"
	"github.com/qdm12/golibs/network"
	"github.com/qdm12/private-internet-access-docker/internal/params"
	"github.com/qdm12/private-internet-access-docker/internal/publicip"
)

func HealthCheck() error {
	paramsReader := params.NewReader(nil)
	ipStatusFilepath, err := paramsReader.GetIPStatusFilepath()
	if err != nil {
		return err
	}

	// Get all VPN ip addresses from openvpn configuration file
	fileManager := files.NewFileManager()
	b, err := fileManager.ReadFile(string(ipStatusFilepath))
	if err != nil {
		return err
	}
	savedPublicIP := net.ParseIP(string(b))
	publicIP, err := publicip.NewIPGetter(network.NewClient(3 * time.Second)).Get()
	if err != nil {
		return err
	}
	if !publicIP.Equal(savedPublicIP) {
		return fmt.Errorf("Public IP address is %s instead of initial vpn IP address %s", publicIP, savedPublicIP)
	}
	return nil
}
