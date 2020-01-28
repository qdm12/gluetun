package params

import (
	"net"
	"os"

	"github.com/qdm12/golibs/logging"
	libparams "github.com/qdm12/golibs/params"
	"github.com/qdm12/golibs/verification"
	"github.com/qdm12/private-internet-access-docker/internal/constants"
)

type ParamsReader interface {
	GetDNSOverTLS() (DNSOverTLS bool, err error)
	GetDNSMaliciousBlocking() (blocking bool, err error)
	GetDNSSurveillanceBlocking() (blocking bool, err error)
	GetDNSAdsBlocking() (blocking bool, err error)
	GetDNSUnblockedHostnames() (hostnames []string, err error)
	GetExtraSubnets() (extraSubnets []*net.IPNet, err error)
	GetNonRoot() (nonRoot bool, err error)
	GetNetworkProtocol() (protocol constants.NetworkProtocol, err error)
	GetPortForwarding() (activated bool, err error)
	GetPortForwardingStatusFilepath() (filepath string, err error)
	GetPIAEncryption() (constants.PIAEncryption, error)
	GetPIARegion() (constants.PIARegion, error)
	GetShadowSocks() (activated bool, err error)
	GetShadowSocksLog() (activated bool, err error)
	GetShadowSocksPort() (port string, err error)
	GetShadowSocksPassword() (password string, err error)
	GetTinyProxy() (activated bool, err error)
	GetTinyProxyLog() (constants.TinyProxyLogLevel, error)
	GetTinyProxyPort() (port string, err error)
	GetTinyProxyUser() (user string, err error)
	GetTinyProxyPassword() (password string, err error)
	GetUser() (s string, err error)
	GetPassword() (s string, err error)
}

type paramsReader struct {
	envParams libparams.EnvParams
	logger    logging.Logger
	verifier  verification.Verifier
	unsetEnv  func(key string) error
}

func NewParamsReader(logger logging.Logger) ParamsReader {
	return &paramsReader{
		envParams: libparams.NewEnvParams(),
		logger:    logger,
		verifier:  verification.NewVerifier(),
		unsetEnv:  os.Unsetenv,
	}
}
