package settings

import (
	"strings"

	"github.com/qdm12/private-internet-access-docker/internal/constants"
	"github.com/qdm12/private-internet-access-docker/internal/params"
)

// DNS contains settings to configure Unbound for DNS over TLS operation
type DNS struct {
	Enabled           bool
	Provider          constants.DNSProvider
	AllowedHostnames  []string
	BlockMalicious    bool
	BlockSurveillance bool
	BlockAds          bool
}

func (d *DNS) String() string {
	if !d.Enabled {
		return "DNS over TLS settings: disabled"
	}
	blockMalicious, blockSurveillance, blockAds := "disabed", "disabed", "disabed"
	if d.BlockMalicious {
		blockMalicious = "enabled"
	}
	if d.BlockSurveillance {
		blockSurveillance = "enabled"
	}
	if d.BlockAds {
		blockAds = "enabled"
	}
	settingsList := []string{
		"DNS over TLS provider: " + string(d.Provider),
		"Block malicious: " + blockMalicious,
		"Block surveillance: " + blockSurveillance,
		"Block ads: " + blockAds,
		"Allowed hostnames: " + strings.Join(d.AllowedHostnames, ", "),
	}
	return "DNS over TLS settings:\n" + strings.Join(settingsList, "\n |--")
}

// GetDNSSettings obtains DNS over TLS settings from environment variables using the params package.
func GetDNSSettings() (settings DNS, err error) {
	settings.Enabled, err = params.GetDNSOverTLS()
	if err != nil || !settings.Enabled {
		return settings, err
	}
	settings.Provider = constants.DNSProvider("cloudflare") // TODO make variable
	settings.AllowedHostnames, err = params.GetDNSUnblockedHostnames()
	if err != nil {
		return settings, err
	}
	settings.BlockMalicious, err = params.GetDNSMaliciousBlocking()
	if err != nil {
		return settings, err
	}
	settings.BlockSurveillance, err = params.GetDNSSurveillanceBlocking()
	if err != nil {
		return settings, err
	}
	settings.BlockAds, err = params.GetDNSAdsBlocking() // TODO add to list
	if err != nil {
		return settings, err
	}
	return settings, nil
}
