package settings

import (
	"fmt"
	"strings"

	"github.com/qdm12/private-internet-access-docker/internal/models"
	"github.com/qdm12/private-internet-access-docker/internal/params"
)

// DNS contains settings to configure Unbound for DNS over TLS operation
type DNS struct {
	Enabled               bool
	Providers             []models.DNSProvider
	AllowedHostnames      []string
	PrivateAddresses      []string
	BlockMalicious        bool
	BlockSurveillance     bool
	BlockAds              bool
	VerbosityLevel        uint8
	VerbosityDetailsLevel uint8
	ValidationLogLevel    uint8
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
	var providersStr []string
	for _, provider := range d.Providers {
		providersStr = append(providersStr, string(provider))
	}
	settingsList := []string{
		"DNS over TLS settings:",
		"DNS over TLS provider:\n  |--" + strings.Join(providersStr, "\n  |--"),
		"Block malicious: " + blockMalicious,
		"Block surveillance: " + blockSurveillance,
		"Block ads: " + blockAds,
		"Allowed hostnames:\n  |--" + strings.Join(d.AllowedHostnames, "\n  |--"),
		"Private addresses:\n  |--" + strings.Join(d.PrivateAddresses, "\n  |--"),
		"Verbosity level: " + fmt.Sprintf("%d/5", d.VerbosityLevel),
		"Verbosity details level: " + fmt.Sprintf("%d/4", d.VerbosityDetailsLevel),
		"Validation log level: " + fmt.Sprintf("%d/2", d.ValidationLogLevel),
	}
	return strings.Join(settingsList, "\n |--")
}

// GetDNSSettings obtains DNS over TLS settings from environment variables using the params package.
func GetDNSSettings(params params.ParamsReader) (settings DNS, err error) {
	settings.Enabled, err = params.GetDNSOverTLS()
	if err != nil || !settings.Enabled {
		return settings, err
	}
	settings.Providers, err = params.GetDNSOverTLSProviders()
	if err != nil {
		return settings, err
	}
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
	settings.BlockAds, err = params.GetDNSAdsBlocking()
	if err != nil {
		return settings, err
	}
	settings.VerbosityLevel, err = params.GetDNSOverTLSVerbosity()
	if err != nil {
		return settings, err
	}
	settings.VerbosityDetailsLevel, err = params.GetDNSOverTLSVerbosityDetails()
	if err != nil {
		return settings, err
	}
	settings.ValidationLogLevel, err = params.GetDNSOverTLSValidationLogLevel()
	if err != nil {
		return settings, err
	}
	settings.PrivateAddresses = []string{ // TODO make env variable
		"127.0.0.1/8",
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"169.254.0.0/16",
		"::1/128",
		"fc00::/7",
		"fe80::/10",
		"::ffff:0:0/96",
	}
	return settings, nil
}
