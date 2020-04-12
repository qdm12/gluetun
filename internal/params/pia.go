package params

import (
	"fmt"
	"math/rand"

	libparams "github.com/qdm12/golibs/params"
	"github.com/qdm12/private-internet-access-docker/internal/constants"
	"github.com/qdm12/private-internet-access-docker/internal/models"
)

// GetPortForwarding obtains if port forwarding on the VPN provider server
// side is enabled or not from the environment variable PORT_FORWARDING
func (p *reader) GetPortForwarding() (activated bool, err error) {
	s, err := p.envParams.GetEnv("PORT_FORWARDING", libparams.Default("off"))
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
func (p *reader) GetPortForwardingStatusFilepath() (filepath models.Filepath, err error) {
	filepathStr, err := p.envParams.GetPath("PORT_FORWARDING_STATUS_FILE", libparams.Default("/forwarded_port"), libparams.CaseSensitiveValue())
	return models.Filepath(filepathStr), err
}

// GetPIAEncryption obtains the encryption level for the PIA connection
// from the environment variable PIA_ENCRYPTION, and using ENCRYPTION for
// retro compatibility
func (p *reader) GetPIAEncryption() (models.PIAEncryption, error) {
	// Retro-compatibility
	s, err := p.envParams.GetValueIfInside("ENCRYPTION", []string{"normal", "strong", ""})
	if err != nil {
		return "", err
	} else if len(s) != 0 {
		p.logger.Warn("You are using the old environment variable ENCRYPTION, please consider changing it to PIA_ENCRYPTION")
		return models.PIAEncryption(s), nil
	}
	s, err = p.envParams.GetValueIfInside("PIA_ENCRYPTION", []string{"normal", "strong"}, libparams.Default("strong"))
	return models.PIAEncryption(s), err
}

// GetPIARegion obtains the region for the PIA server from the
// environment variable REGION
func (p *reader) GetPIARegion() (region models.PIARegion, err error) {
	choices := append(constants.PIAGeoChoices(), "")
	s, err := p.envParams.GetValueIfInside("REGION", choices)
	if len(s) == 0 { // Suggestion by @rorph https://github.com/rorph
		s = choices[rand.Int()%len(choices)] //nolint:gosec
	}
	return models.PIARegion(s), err
}
