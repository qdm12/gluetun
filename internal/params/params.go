package params

import (
	"net"
	"os"
	"time"

	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/golibs/files"
	"github.com/qdm12/golibs/logging"
	libparams "github.com/qdm12/golibs/params"
	"github.com/qdm12/golibs/verification"
)

// Reader contains methods to obtain parameters.
type Reader interface {
	GetVPNSP() (vpnServiceProvider models.VPNProvider, err error)

	// DNS over TLS getters
	GetDNSOverTLS() (DNSOverTLS bool, err error)
	GetDNSOverTLSProviders() (providers []models.DNSProvider, err error)
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
	GetUID() (uid int, err error)
	GetGID() (gid int, err error)
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
	GetPassword(required bool) (s string, err error)
	GetNetworkProtocol() (protocol models.NetworkProtocol, err error)
	GetOpenVPNVerbosity() (verbosity int, err error)
	GetOpenVPNRoot() (root bool, err error)
	GetTargetIP() (ip net.IP, err error)
	GetOpenVPNCipher() (cipher string, err error)
	GetOpenVPNAuth() (auth string, err error)
	GetOpenVPNIPv6() (tunnel bool, err error)

	// PIA getters
	GetPortForwarding() (activated bool, err error)
	GetPortForwardingStatusFilepath() (filepath models.Filepath, err error)
	GetPIAEncryptionPreset() (preset string, err error)
	GetPIARegions() (regions []string, err error)

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
	GetShadowSocksPort() (port uint16, err error)
	GetShadowSocksPassword() (password string, err error)
	GetShadowSocksMethod() (method string, err error)

	// HTTP proxy getters
	GetHTTPProxy() (activated bool, err error)
	GetHTTPProxyLog() (log bool, err error)
	GetHTTPProxyPort() (port uint16, err error)
	GetHTTPProxyUser() (user string, err error)
	GetHTTPProxyPassword() (password string, err error)
	GetHTTPProxyStealth() (stealth bool, err error)

	// Public IP getters
	GetPublicIPPeriod() (period time.Duration, err error)

	// Control server
	GetControlServerPort() (port uint16, err error)
	GetControlServerLog() (enabled bool, err error)

	GetVersionInformation() (enabled bool, err error)

	GetUpdaterPeriod() (period time.Duration, err error)
}

type reader struct {
	envParams   libparams.EnvParams
	logger      logging.Logger
	verifier    verification.Verifier
	unsetEnv    func(key string) error
	fileManager files.FileManager
}

// Newreader returns a paramsReadeer object to read parameters from
// environment variables.
func NewReader(logger logging.Logger, fileManager files.FileManager) Reader {
	return &reader{
		envParams:   libparams.NewEnvParams(),
		logger:      logger,
		verifier:    verification.NewVerifier(),
		unsetEnv:    os.Unsetenv,
		fileManager: fileManager,
	}
}

// GetVPNSP obtains the VPN service provider to use from the environment variable VPNSP.
func (r *reader) GetVPNSP() (vpnServiceProvider models.VPNProvider, err error) {
	s, err := r.envParams.GetValueIfInside(
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
	return r.envParams.GetOnOff("VERSION_INFORMATION", libparams.Default("on"))
}

func (r *reader) onRetroActive(oldKey, newKey string) {
	r.logger.Warn(
		"You are using the old environment variable %s, please consider changing it to %s",
		oldKey, newKey,
	)
}
