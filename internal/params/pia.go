package params

import (
	"fmt"

	libparams "github.com/qdm12/golibs/params"
	"github.com/qdm12/private-internet-access-docker/internal/constants"
)

// GetPortForwarding obtains if port forwarding on the VPN provider server
// side is enabled or not from the environment variable PORT_FORWARDING
func GetPortForwarding(envParams libparams.EnvParams) (activated bool, err error) {
	s, err := envParams.GetEnv("PORT_FORWARDING", libparams.Default("off"))
	if err != nil {
		return false, err
	}
	// Custom for retro-compatibility
	if s == "false" || s == "off" {
		return false, nil
	} else if s == "true" || s == "on" {
		return true, nil
	}
	return false, fmt.Errorf("PORT_FORWARDING can only be \"on\" or \"off\"")
}

// GetPortForwardingStatusFilepath obtains the port forwarding status file path
// from the environment variable PORT_FORWARDING_STATUS_FILE
func GetPortForwardingStatusFilepath(envParams libparams.EnvParams) (filepath string, err error) {
	return envParams.GetPath("PORT_FORWARDING_STATUS_FILE", libparams.Default("/forwarded_port"))
}

// GetPIAEncryption obtains the encryption level for the PIA connection
// from the environment variable ENCRYPTION
func GetPIAEncryption(envParams libparams.EnvParams) (constants.PIAEncryption, error) {
	s, err := envParams.GetValueIfInside("ENCRYPTION", []string{"normal", "strong"}, libparams.Default("strong"))
	return constants.PIAEncryption(s), err
}

// GetPIARegion obtains the region for the PIA server from the
// environment variable REGION
func GetPIARegion(envParams libparams.EnvParams) (constants.PIARegion, error) {
	s, err := envParams.GetValueIfInside("REGION", []string{"Netherlands"}, libparams.Compulsory())
	return constants.PIARegion(s), err
}
