package cli

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/qdm12/golibs/files"
	"github.com/qdm12/golibs/network"
	"github.com/qdm12/private-internet-access-docker/internal/params"
)

func HealthCheck() error {
	paramsReader := params.NewReader(nil)
	ipStatusFilepath, err := paramsReader.GetIPStatusFilepath()
	if err != nil {
		return err
	}
	// Get VPN ip address written to file
	fileManager := files.NewFileManager()
	b, err := fileManager.ReadFile(string(ipStatusFilepath))
	if err != nil {
		return err
	}
	vpnIP := string(b)

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
	if publicIP != vpnIP {
		return fmt.Errorf("Public IP address %s does not match VPN ip address %s on file", publicIP, vpnIP)
	}
	return nil
}
