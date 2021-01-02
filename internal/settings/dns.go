package settings

import (
	"fmt"
	"net"
	"strings"
	"time"

	unboundmodels "github.com/qdm12/dns/pkg/models"
	unbound "github.com/qdm12/dns/pkg/unbound"
	"github.com/qdm12/gluetun/internal/params"
)

// DNS contains settings to configure Unbound for DNS over TLS operation.
type DNS struct { //nolint:maligned
	Enabled           bool
	PlaintextAddress  net.IP
	KeepNameserver    bool
	BlockMalicious    bool
	BlockAds          bool
	BlockSurveillance bool
	UpdatePeriod      time.Duration
	Unbound           unboundmodels.Settings
}

func (d *DNS) String() string {
	if !d.Enabled {
		return fmt.Sprintf("DNS over TLS disabled, using plaintext DNS %s", d.PlaintextAddress)
	}
	blockMalicious, blockSurveillance, blockAds := disabled, disabled, disabled
	if d.BlockMalicious {
		blockMalicious = enabled
	}
	if d.BlockSurveillance {
		blockSurveillance = enabled
	}
	if d.BlockAds {
		blockAds = enabled
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
		"DNS settings:",
		"Block malicious: " + blockMalicious,
		"Block surveillance: " + blockSurveillance,
		"Block ads: " + blockAds,
		"Update: " + update,
		"Keep nameserver (disabled blocking): " + keepNameserver,
		"Unbound settings: " + "\n   |--" + strings.Join(d.Unbound.Lines(), "\n   |--"),
	}
	return strings.Join(settingsList, "\n |--")
}

// GetDNSSettings obtains DNS over TLS settings from environment variables using the params package.
func GetDNSSettings(paramsReader params.Reader) (settings DNS, err error) {
	settings.Enabled, err = paramsReader.GetDNSOverTLS()
	if err != nil {
		return settings, err
	}

	// Plain DNS settings
	settings.PlaintextAddress, err = paramsReader.GetDNSPlaintext()
	if err != nil {
		return settings, err
	}
	settings.KeepNameserver, err = paramsReader.GetDNSKeepNameserver()
	if err != nil {
		return settings, err
	}

	// DNS over TLS external settings
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
	settings.UpdatePeriod, err = paramsReader.GetDNSUpdatePeriod()
	if err != nil {
		return settings, err
	}

	// Unbound specific settings
	settings.Unbound, err = getUnboundSettings(paramsReader)
	if err != nil {
		return settings, err
	}

	// Consistency check
	IPv6Support := false
	for _, provider := range settings.Unbound.Providers {
		providerData, ok := unbound.GetProviderData(provider)
		switch {
		case !ok:
			return settings, fmt.Errorf("DNS provider %q does not have associated data", provider)
		case !providerData.SupportsTLS:
			return settings, fmt.Errorf("DNS provider %q does not support DNS over TLS", provider)
		case providerData.SupportsIPv6:
			IPv6Support = true
		}
	}
	if settings.Unbound.IPv6 && !IPv6Support {
		return settings, fmt.Errorf("None of the DNS over TLS provider(s) set support IPv6")
	}
	return settings, nil
}

func getUnboundSettings(reader params.Reader) (settings unboundmodels.Settings, err error) {
	settings.Providers, err = reader.GetDNSOverTLSProviders()
	if err != nil {
		return settings, err
	}
	settings.ListeningPort = 53
	settings.Caching, err = reader.GetDNSOverTLSCaching()
	if err != nil {
		return settings, err
	}
	settings.IPv4 = true
	settings.IPv6, err = reader.GetDNSOverTLSIPv6()
	if err != nil {
		return settings, err
	}
	settings.VerbosityLevel, err = reader.GetDNSOverTLSVerbosity()
	if err != nil {
		return settings, err
	}
	settings.VerbosityDetailsLevel, err = reader.GetDNSOverTLSVerbosityDetails()
	if err != nil {
		return settings, err
	}
	settings.ValidationLogLevel, err = reader.GetDNSOverTLSValidationLogLevel()
	if err != nil {
		return settings, err
	}
	settings.BlockedHostnames = []string{}
	settings.BlockedIPs, err = reader.GetDNSOverTLSPrivateAddresses()
	if err != nil {
		return settings, err
	}
	settings.AllowedHostnames, err = reader.GetDNSUnblockedHostnames()
	if err != nil {
		return settings, err
	}
	return settings, nil
}
