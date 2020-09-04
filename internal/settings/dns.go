package settings

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/params"
)

// DNS contains settings to configure Unbound for DNS over TLS operation
type DNS struct {
	Enabled               bool
	KeepNameserver        bool
	Providers             []models.DNSProvider
	PlaintextAddress      net.IP
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
	UpdatePeriod          time.Duration
}

func (d *DNS) String() string {
	if !d.Enabled {
		return fmt.Sprintf("DNS over TLS disabled, using plaintext DNS %s", d.PlaintextAddress)
	}
	caching, blockMalicious, blockSurveillance, blockAds, ipv6 := disabled, disabled, disabled, disabled, disabled
	if d.Caching {
		caching = enabled
	}
	if d.BlockMalicious {
		blockMalicious = enabled
	}
	if d.BlockSurveillance {
		blockSurveillance = enabled
	}
	if d.BlockAds {
		blockAds = enabled
	}
	if d.IPv6 {
		ipv6 = enabled
	}
	providersStr := make([]string, len(d.Providers))
	for i := range d.Providers {
		providersStr[i] = string(d.Providers[i])
	}
	update := "deactivated"
	if d.UpdatePeriod > 0 {
		update = fmt.Sprintf("every %s", d.UpdatePeriod)
	}
	keepNameserver := "no"
	if d.KeepNameserver {
		keepNameserver = "yes"
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
		"Update: " + update,
		"Keep nameserver (disabled blocking): " + keepNameserver,
	}
	return strings.Join(settingsList, "\n |--")
}

// GetDNSSettings obtains DNS over TLS settings from environment variables using the params package.
func GetDNSSettings(paramsReader params.Reader) (settings DNS, err error) {
	settings.Enabled, err = paramsReader.GetDNSOverTLS()
	if err != nil {
		return settings, err
	}
	if !settings.Enabled {
		settings.PlaintextAddress, err = paramsReader.GetDNSPlaintext()
		return settings, err
	}
	settings.Providers, err = paramsReader.GetDNSOverTLSProviders()
	if err != nil {
		return settings, err
	}
	settings.AllowedHostnames, err = paramsReader.GetDNSUnblockedHostnames()
	if err != nil {
		return settings, err
	}
	settings.Caching, err = paramsReader.GetDNSOverTLSCaching()
	if err != nil {
		return settings, err
	}
	settings.BlockMalicious, err = paramsReader.GetDNSMaliciousBlocking()
	if err != nil {
		return settings, err
	}
	settings.BlockSurveillance, err = paramsReader.GetDNSSurveillanceBlocking()
	if err != nil {
		return settings, err
	}
	settings.BlockAds, err = paramsReader.GetDNSAdsBlocking()
	if err != nil {
		return settings, err
	}
	settings.VerbosityLevel, err = paramsReader.GetDNSOverTLSVerbosity()
	if err != nil {
		return settings, err
	}
	settings.VerbosityDetailsLevel, err = paramsReader.GetDNSOverTLSVerbosityDetails()
	if err != nil {
		return settings, err
	}
	settings.ValidationLogLevel, err = paramsReader.GetDNSOverTLSValidationLogLevel()
	if err != nil {
		return settings, err
	}
	settings.PrivateAddresses, err = paramsReader.GetDNSOverTLSPrivateAddresses()
	if err != nil {
		return settings, err
	}
	settings.IPv6, err = paramsReader.GetDNSOverTLSIPv6()
	if err != nil {
		return settings, err
	}
	settings.UpdatePeriod, err = paramsReader.GetDNSUpdatePeriod()
	if err != nil {
		return settings, err
	}
	settings.KeepNameserver, err = paramsReader.GetDNSKeepNameserver()
	if err != nil {
		return settings, err
	}

	// Consistency check
	IPv6Support := false
	for _, provider := range settings.Providers {
		providerData, ok := constants.DNSProviderMapping()[provider]
		switch {
		case !ok:
			return settings, fmt.Errorf("DNS provider %q does not have associated data", provider)
		case !providerData.SupportsTLS:
			return settings, fmt.Errorf("DNS provider %q does not support DNS over TLS", provider)
		case providerData.SupportsIPv6:
			IPv6Support = true
		}
	}
	if settings.IPv6 && !IPv6Support {
		return settings, fmt.Errorf("None of the DNS over TLS provider(s) set support IPv6")
	}
	return settings, nil
}
