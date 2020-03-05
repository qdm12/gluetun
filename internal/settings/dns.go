package settings

import (
	"fmt"
	"strings"

	"github.com/qdm12/private-internet-access-docker/internal/constants"
	"github.com/qdm12/private-internet-access-docker/internal/models"
	"github.com/qdm12/private-internet-access-docker/internal/params"
)

// DNS contains settings to configure Unbound for DNS over TLS operation
type DNS struct {
	Enabled               bool
	Providers             []models.DNSProvider
	AllowedHostnames      []string
	PrivateAddresses      []string
	Caching               bool
	BlockMalicious        bool
	BlockSurveillance     bool
	BlockAds              bool
	VerbosityLevel        uint8
	VerbosityDetailsLevel uint8
	ValidationLogLevel    uint8
	IPv6                  bool
}

func (d *DNS) String() string {
	if !d.Enabled {
		return "DNS over TLS settings: disabled"
	}
	caching, blockMalicious, blockSurveillance, blockAds, ipv6 := "disabled", "disabed", "disabed", "disabed", "disabed"
	if d.Caching {
		caching = "enabled"
	}
	if d.BlockMalicious {
		blockMalicious = "enabled"
	}
	if d.BlockSurveillance {
		blockSurveillance = "enabled"
	}
	if d.BlockAds {
		blockAds = "enabled"
	}
	if d.IPv6 {
		ipv6 = "enabled"
	}
	var providersStr []string
	for _, provider := range d.Providers {
		providersStr = append(providersStr, string(provider))
	}
	settingsList := []string{
		"DNS over TLS settings:",
		"DNS over TLS provider:\n  |--" + strings.Join(providersStr, "\n  |--"),
		"Caching: " + caching,
		"Block malicious: " + blockMalicious,
		"Block surveillance: " + blockSurveillance,
		"Block ads: " + blockAds,
		"Allowed hostnames:\n  |--" + strings.Join(d.AllowedHostnames, "\n  |--"),
		"Private addresses:\n  |--" + strings.Join(d.PrivateAddresses, "\n  |--"),
		"Verbosity level: " + fmt.Sprintf("%d/5", d.VerbosityLevel),
		"Verbosity details level: " + fmt.Sprintf("%d/4", d.VerbosityDetailsLevel),
		"Validation log level: " + fmt.Sprintf("%d/2", d.ValidationLogLevel),
		"IPv6 resolution: " + ipv6,
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
	settings.Caching, err = params.GetDNSOverTLSCaching()
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
	settings.PrivateAddresses = params.GetDNSOverTLSPrivateAddresses()
	settings.IPv6, err = params.GetDNSOverTLSIPv6()
	if err != nil {
		return settings, err
	}

	// Consistency check
	IPv6Support := false
	for _, provider := range settings.Providers {
		providerData, ok := constants.DNSProviderMapping()[provider]
		if !ok {
			return settings, fmt.Errorf("DNS provider %q does not have associated data", provider)
		} else if !providerData.SupportsTLS {
			return settings, fmt.Errorf("DNS provider %q does not support DNS over TLS", provider)
		} else if providerData.SupportsIPv6 {
			IPv6Support = true
		}
	}
	if settings.IPv6 && !IPv6Support {
		return settings, fmt.Errorf("None of the DNS over TLS provider(s) set support IPv6")
	}
	return settings, nil
}
