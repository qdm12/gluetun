package params

import libparams "github.com/qdm12/golibs/params"

import "strings"

// GetDNSOverTLS obtains if the DNS over TLS should be enabled
// from the environment variable DOT
func GetDNSOverTLS() (DNSOverTLS bool, err error) {
	return libparams.GetOnOff("DOT", true)
}

// GetDNSMaliciousBlocking obtains if malicious hostnames/IPs should be blocked
// from being resolved by Unbound, using the environment variable BLOCK_MALICIOUS
func GetDNSMaliciousBlocking() (blocking bool, err error) {
	return libparams.GetOnOff("BLOCK_MALICIOUS", false)
}

// GetDNSSurveillanceBlocking obtainsv if surveillance hostnames/IPs should be blocked
// from being resolved by Unbound, using the environment variable BLOCK_NSA
func GetDNSSurveillanceBlocking() (blocking bool, err error) {
	return libparams.GetOnOff("BLOCK_NSA", false)
}

// GetUnblockedHostnames obtains a list of hostnames to unblock from block lists
// from the comma separated list for the environment variable UNBLOCK
func GetUnblockedHostnames() (hostnames []string, err error) {
	s := libparams.GetEnv("UNBLOCK", "")
	if len(s) == 0 {
		return nil, nil
	}
	hostnames = strings.Split(s, ",")
	// TODO validate hostnames
	return hostnames, nil
}
