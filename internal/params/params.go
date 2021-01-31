package params

import (
	"net"
	"time"

	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/os"
	libparams "github.com/qdm12/golibs/params"
	"github.com/qdm12/golibs/verification"
)

// Reader contains methods to obtain parameters.
type Reader interface {
	GetVPNSP() (vpnServiceProvider models.VPNProvider, err error)

	// DNS over TLS getters
	GetDNSOverTLS() (DNSOverTLS bool, err error)
	GetDNSOverTLSProviders() (providers []string, err error)
	GetDNSOverTLSCaching() (caching bool, err error)
	GetDNSOverTLSVerbosity() (verbosityLevel uint8, err error)
	GetDNSOverTLSVerbosityDetails() (verbosityDetailsLevel uint8, err error)
	GetDNSOverTLSValidationLogLevel() (validationLogLevel uint8, err error)
	GetDNSMaliciousBlocking() (blocking bool, err error)
	GetDNSSurveillanceBlocking() (blocking bool, err error)
	GetDNSAdsBlocking() (blocking bool, err error)
	GetDNSUnblockedHostnames() (hostnames []string, err error)
	GetDNSOverTLSPrivateAddresses() (privateAddresses []string, err error)
	GetDNSOverTLSIPv6() (ipv6 bool, err error)
	GetDNSUpdatePeriod() (period time.Duration, err error)
	GetDNSPlaintext() (ip net.IP, err error)
	GetDNSKeepNameserver() (on bool, err error)

	// System
	GetPUID() (puid int, err error)
	GetPGID() (pgid int, err error)
	GetTimezone() (timezone string, err error)
	GetPublicIPFilepath() (filepath models.Filepath, err error)

	// Firewall getters
	GetFirewall() (enabled bool, err error)
	GetVPNInputPorts() (ports []uint16, err error)
	GetInputPorts() (ports []uint16, err error)
	GetOutboundSubnets() (outboundSubnets []net.IPNet, err error)
	GetFirewallDebug() (debug bool, err error)

	// VPN getters
	GetUser() (s string, err error)
	GetPassword() (s string, err error)
	GetNetworkProtocol() (protocol models.NetworkProtocol, err error)
	GetOpenVPNVerbosity() (verbosity int, err error)
	GetOpenVPNRoot() (root bool, err error)
	GetTargetIP() (ip net.IP, err error)
	GetOpenVPNCipher() (cipher string, err error)
	GetOpenVPNAuth() (auth string, err error)
	GetOpenVPNIPv6() (tunnel bool, err error)
	GetOpenVPNMSSFix() (mssFix uint16, err error)

	// PIA getters
	GetPortForwarding() (activated bool, err error)
	GetPortForwardingStatusFilepath() (filepath models.Filepath, err error)
	GetPIAEncryptionPreset() (preset string, err error)
	GetPIARegions() (regions []string, err error)
	GetPIAPort() (port uint16, err error)

	// Mullvad getters
	GetMullvadCountries() (countries []string, err error)
	GetMullvadCities() (cities []string, err error)
	GetMullvadISPs() (ips []string, err error)
	GetMullvadPort() (port uint16, err error)
	GetMullvadOwned() (owned bool, err error)

	// Windscribe getters
	GetWindscribeRegions() (countries []string, err error)
	GetWindscribeCities() (cities []string, err error)
	GetWindscribeHostnames() (hostnames []string, err error)
	GetWindscribePort(protocol models.NetworkProtocol) (port uint16, err error)

	// Surfshark getters
	GetSurfsharkRegions() (countries []string, err error)

	// Cyberghost getters
	GetCyberghostGroup() (group string, err error)
	GetCyberghostRegions() (regions []string, err error)
	GetCyberghostClientKey() (clientKey string, err error)
	GetCyberghostClientCertificate() (clientCertificate string, err error)

	// Vyprvpn getters
	GetVyprvpnRegions() (regions []string, err error)

	// NordVPN getters
	GetNordvpnRegions() (regions []string, err error)
	GetNordvpnNumbers() (numbers []uint16, err error)

	// Privado getters
	GetPrivadoHostnames() (hostnames []string, err error)

	// PureVPN getters
	GetPurevpnRegions() (regions []string, err error)
	GetPurevpnCountries() (countries []string, err error)
	GetPurevpnCities() (cities []string, err error)

	// Shadowsocks getters
	GetShadowSocks() (activated bool, err error)
	GetShadowSocksLog() (activated bool, err error)
	GetShadowSocksPort() (port uint16, warning string, err error)
	GetShadowSocksPassword() (password string, err error)
	GetShadowSocksMethod() (method string, err error)

	// HTTP proxy getters
	GetHTTPProxy() (activated bool, err error)
	GetHTTPProxyLog() (log bool, err error)
	GetHTTPProxyPort() (port uint16, warning string, err error)
	GetHTTPProxyUser() (user string, err error)
	GetHTTPProxyPassword() (password string, err error)
	GetHTTPProxyStealth() (stealth bool, err error)

	// Public IP getters
	GetPublicIPPeriod() (period time.Duration, err error)

	// Control server
	GetControlServerPort() (port uint16, warning string, err error)
	GetControlServerLog() (enabled bool, err error)

	GetVersionInformation() (enabled bool, err error)

	GetUpdaterPeriod() (period time.Duration, err error)
}

type reader struct {
	env    libparams.Env
	logger logging.Logger
	regex  verification.Regex
	os     os.OS
}

// Newreader returns a paramsReadeer object to read parameters from
// environment variables.
func NewReader(logger logging.Logger, os os.OS) Reader {
	return &reader{
		env:    libparams.NewEnv(),
		logger: logger,
		regex:  verification.NewRegex(),
		os:     os,
	}
}

// GetVPNSP obtains the VPN service provider to use from the environment variable VPNSP.
func (r *reader) GetVPNSP() (vpnServiceProvider models.VPNProvider, err error) {
	s, err := r.env.Inside(
		"VPNSP",
		[]string{
			"pia", "private internet access",
			"mullvad", "windscribe", "surfshark", "cyberghost",
			"vyprvpn", "nordvpn", "purevpn", "privado",
		}, libparams.Default("private internet access"))
	if s == "pia" {
		s = "private internet access"
	}
	return models.VPNProvider(s), err
}

func (r *reader) GetVersionInformation() (enabled bool, err error) {
	return r.env.OnOff("VERSION_INFORMATION", libparams.Default("on"))
}

func (r *reader) onRetroActive(oldKey, newKey string) {
	r.logger.Warn(
		"You are using the old environment variable %s, please consider changing it to %s",
		oldKey, newKey,
	)
}
