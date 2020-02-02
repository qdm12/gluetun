package params

import (
	"fmt"

	libparams "github.com/qdm12/golibs/params"
	"github.com/qdm12/private-internet-access-docker/internal/constants"
	"github.com/qdm12/private-internet-access-docker/internal/models"
)

// GetUser obtains the user to use to connect to the VPN servers
func (p *paramsReader) GetUser() (s string, err error) {
	defer func() {
		unsetenvErr := p.unsetEnv("USER")
		if err == nil {
			err = unsetenvErr
		}
	}()
	s, err = p.envParams.GetEnv("USER")
	if err != nil {
		return "", err
	} else if len(s) == 0 {
		return s, fmt.Errorf("USER environment variable cannot be empty")
	}
	return s, nil
}

// GetPassword obtains the password to use to connect to the VPN servers
func (p *paramsReader) GetPassword() (s string, err error) {
	defer func() {
		unsetenvErr := p.unsetEnv("PASSWORD")
		if err == nil {
			err = unsetenvErr
		}
	}()
	s, err = p.envParams.GetEnv("PASSWORD")
	if err != nil {
		return "", err
	} else if len(s) == 0 {
		return s, fmt.Errorf("PASSWORD environment variable cannot be empty")
	}
	return s, nil
}

// GetPortForwarding obtains if port forwarding on the VPN provider server
// side is enabled or not from the environment variable PORT_FORWARDING
func (p *paramsReader) GetPortForwarding() (activated bool, err error) {
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
func (p *paramsReader) GetPortForwardingStatusFilepath() (filepath models.Filepath, err error) {
	filepathStr, err := p.envParams.GetPath("PORT_FORWARDING_STATUS_FILE", libparams.Default("/forwarded_port"))
	return models.Filepath(filepathStr), err
}

// GetPIAEncryption obtains the encryption level for the PIA connection
// from the environment variable ENCRYPTION
func (p *paramsReader) GetPIAEncryption() (models.PIAEncryption, error) {
	s, err := p.envParams.GetValueIfInside("ENCRYPTION", []string{"normal", "strong"}, libparams.Default("strong"))
	return models.PIAEncryption(s), err
}

// GetPIARegion obtains the region for the PIA server from the
// environment variable REGION
func (p *paramsReader) GetPIARegion() (region models.PIARegion, err error) {
	s, err := p.envParams.GetValueIfInside("REGION", []string{
		string(constants.AUMelbourne), string(constants.AUPerth), string(constants.AUSydney), string(constants.Austria), string(constants.Belgium), string(constants.CAMontreal), string(constants.CAToronto), string(constants.CAVancouver), string(constants.CzechRepublic), string(constants.DEBerlin), string(constants.DEFrankfurt), string(constants.Denmark), string(constants.Finland), string(constants.France), string(constants.HongKong), string(constants.Hungary), string(constants.India), string(constants.Ireland), string(constants.Israel), string(constants.Italy), string(constants.Japan), string(constants.Luxembourg), string(constants.Mexico), string(constants.Netherlands), string(constants.NewZealand), string(constants.Norway), string(constants.Poland), string(constants.Romania), string(constants.Singapore), string(constants.Spain), string(constants.Sweden), string(constants.Switzerland), string(constants.UAE), string(constants.UKLondon), string(constants.UKManchester), string(constants.UKSouthampton), string(constants.USAtlanta), string(constants.USCalifornia), string(constants.USChicago), string(constants.USDenver), string(constants.USEast), string(constants.USFlorida), string(constants.USHouston), string(constants.USLasVegas), string(constants.USNewYorkCity), string(constants.USSeattle), string(constants.USSiliconValley), string(constants.USTexas), string(constants.USWashingtonDC), string(constants.USWest),
	}, libparams.Compulsory())
	return models.PIARegion(s), err
}
