package params

import (
	"net"
	"os"
	"time"

	"github.com/qdm12/golibs/logging"
	libparams "github.com/qdm12/golibs/params"
	"github.com/qdm12/golibs/verification"
	"github.com/qdm12/private-internet-access-docker/internal/models"
)

// Reader contains methods to obtain parameters
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

	// System
	GetUID() (uid int, err error)
	GetGID() (gid int, err error)
	GetTimezone() (timezone string, err error)
	GetIPStatusFilepath() (filepath models.Filepath, err error)

	// Firewall getters
	GetExtraSubnets() (extraSubnets []net.IPNet, err error)

	// VPN getters
	GetUser() (s string, err error)
	GetPassword() (s string, err error)
	GetNetworkProtocol() (protocol models.NetworkProtocol, err error)
	GetOpenVPNVerbosity() (verbosity int, err error)
	GetOpenVPNRoot() (root bool, err error)
	GetTargetIP() (ip net.IP, err error)
	GetOpenVPNCipher() (cipher string, err error)
	GetOpenVPNAuth() (auth string, err error)

	// PIA getters
	GetPortForwarding() (activated bool, err error)
	GetPortForwardingStatusFilepath() (filepath models.Filepath, err error)
	GetPIAEncryption() (models.PIAEncryption, error)
	GetPIARegion() (models.PIARegion, error)

	// Mullvad getters
	GetMullvadCountry() (country models.MullvadCountry, err error)
	GetMullvadCity() (country models.MullvadCity, err error)
	GetMullvadISP() (country models.MullvadProvider, err error)
	GetMullvadPort() (port uint16, err error)

	// Windscribe getters
	GetWindscribeRegion() (country models.WindscribeRegion, err error)
	GetWindscribePort(protocol models.NetworkProtocol) (port uint16, err error)

	// Shadowsocks getters
	GetShadowSocks() (activated bool, err error)
	GetShadowSocksLog() (activated bool, err error)
	GetShadowSocksPort() (port uint16, err error)
	GetShadowSocksPassword() (password string, err error)
	GetShadowSocksMethod() (method string, err error)

	// Tinyproxy getters
	GetTinyProxy() (activated bool, err error)
	GetTinyProxyLog() (models.TinyProxyLogLevel, error)
	GetTinyProxyPort() (port uint16, err error)
	GetTinyProxyUser() (user string, err error)
	GetTinyProxyPassword() (password string, err error)

	// Version getters
	GetVersion() string
	GetBuildDate() string
	GetVcsRef() string
}

type reader struct {
	envParams libparams.EnvParams
	logger    logging.Logger
	verifier  verification.Verifier
	unsetEnv  func(key string) error
}

// Newreader returns a paramsReadeer object to read parameters from
// environment variables
func NewReader(logger logging.Logger) Reader {
	return &reader{
		envParams: libparams.NewEnvParams(),
		logger:    logger,
		verifier:  verification.NewVerifier(),
		unsetEnv:  os.Unsetenv,
	}
}

// GetVPNSP obtains the VPN service provider to use from the environment variable VPNSP
func (p *reader) GetVPNSP() (vpnServiceProvider models.VPNProvider, err error) {
	s, err := p.envParams.GetValueIfInside("VPNSP", []string{"pia", "private internet access", "mullvad", "windscribe"})
	if s == "pia" {
		s = "private internet access"
	}
	return models.VPNProvider(s), err
}
