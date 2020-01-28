package params

import libparams "github.com/qdm12/golibs/params"

import "strings"

// GetDNSOverTLS obtains if the DNS over TLS should be enabled
// from the environment variable DOT
func (p *paramsReader) GetDNSOverTLS() (DNSOverTLS bool, err error) {
	return p.envParams.GetOnOff("DOT", libparams.Default("on"))
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
	// TODO validate hostnames
	return hostnames, nil
}
