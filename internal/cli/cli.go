package cli

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/qdm12/golibs/files"
	"github.com/qdm12/golibs/network"
	"github.com/qdm12/private-internet-access-docker/internal/constants"
)

func HealthCheck() error {
	// Get all VPN ip addresses from openvpn configuration file
	fileManager := files.NewFileManager()
	b, err := fileManager.ReadFile(string(constants.OpenVPNConf))
	if err != nil {
		return err
	}
	var vpnIPs []string
	for _, line := range strings.Split(string(b), "\n") {
		if strings.HasPrefix(line, "remote ") {
			fields := strings.Fields(line)
			vpnIPs = append(vpnIPs, fields[1])
		}
	}

	// Get public IP address from one of the following urls
	urls := []string{
		"http://ip1.dynupdate.no-ip.com:8245",
		"http://ip1.dynupdate.no-ip.com",
		"https://api.ipify.org",
		"https://diagnostic.opendns.com/myip",
		"https://domains.google.com/checkip",
		"https://ifconfig.io/ip",
		"https://ip4.ddnss.de/meineip.php",
		"https://ipinfo.io/ip",
	}
	url := urls[rand.Intn(len(urls))]
	client := network.NewClient(3 * time.Second)
	content, status, err := client.GetContent(url, network.UseRandomUserAgent())
	if err != nil {
		return err
	} else if status != http.StatusOK {
		return fmt.Errorf("Received unexpected status code %d from %s", status, url)
	}
	publicIP := strings.ReplaceAll(string(content), "\n", "")
	match := false
	for _, vpnIP := range vpnIPs {
		if publicIP == vpnIP {
			match = true
			break
		}
	}
	if !match {
		return fmt.Errorf("Public IP address %s does not match any of the VPN ip addresses %s", publicIP, strings.Join(vpnIPs, ", "))
	}
	return nil
}
