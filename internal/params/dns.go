package params

import (
	"fmt"
	"strings"

	libparams "github.com/qdm12/golibs/params"
	"github.com/qdm12/private-internet-access-docker/internal/constants"
	"github.com/qdm12/private-internet-access-docker/internal/models"
)

// GetDNSOverTLS obtains if the DNS over TLS should be enabled
// from the environment variable DOT
func (p *paramsReader) GetDNSOverTLS() (DNSOverTLS bool, err error) {
	return p.envParams.GetOnOff("DOT", libparams.Default("on"))
}

// GetDNSOverTLSProviders obtains the DNS over TLS providers to use
// from the environment variable DOT_PROVIDERS
func (p *paramsReader) GetDNSOverTLSProviders() (providers []models.DNSProvider, err error) {
	s, err := p.envParams.GetEnv("DOT_PROVIDERS", libparams.Default("cloudflare"))
	if err != nil {
		return nil, err
	}
	for _, word := range strings.Split(s, ",") {
		provider := models.DNSProvider(word)
		switch provider {
		case constants.Cloudflare, constants.Google, constants.Quad9, constants.Quadrant, constants.CleanBrowsing, constants.SecureDNS, constants.LibreDNS:
			providers = append(providers, provider)
		default:
			return nil, fmt.Errorf("DNS over TLS provider %q is not valid", provider)
		}
	}
	return providers, nil
}

// GetDNSMaliciousBlocking obtains if malicious hostnames/IPs should be blocked
// from being resolved by Unbound, using the environment variable BLOCK_MALICIOUS
func (p *paramsReader) GetDNSMaliciousBlocking() (blocking bool, err error) {
	return p.envParams.GetOnOff("BLOCK_MALICIOUS", libparams.Default("off"))
}

// GetDNSSurveillanceBlocking obtains if surveillance hostnames/IPs should be blocked
// from being resolved by Unbound, using the environment variable BLOCK_NSA
func (p *paramsReader) GetDNSSurveillanceBlocking() (blocking bool, err error) {
	return p.envParams.GetOnOff("BLOCK_NSA", libparams.Default("off"))
}

// GetDNSAdsBlocking obtains if ads hostnames/IPs should be blocked
// from being resolved by Unbound, using the environment variable BLOCK_ADS
func (p *paramsReader) GetDNSAdsBlocking() (blocking bool, err error) {
	return p.envParams.GetOnOff("BLOCK_ADS", libparams.Default("off"))
}

// GetDNSUnblockedHostnames obtains a list of hostnames to unblock from block lists
// from the comma separated list for the environment variable UNBLOCK
func (p *paramsReader) GetDNSUnblockedHostnames() (hostnames []string, err error) {
	s, err := p.envParams.GetEnv("UNBLOCK")
	if err != nil {
		return nil, err
	}
	if len(s) == 0 {
		return nil, nil
	}
	hostnames = strings.Split(s, ",")
	for _, hostname := range hostnames {
		if !p.verifier.MatchHostname(hostname) {
			return nil, fmt.Errorf("hostname %q does not seem valid", hostname)
		}
	}
	return hostnames, nil
}
