package params

import (
	"fmt"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	libparams "github.com/qdm12/golibs/params"
)

// GetPortForwarding obtains if port forwarding on the VPN provider server
// side is enabled or not from the environment variable PORT_FORWARDING
// Only valid for older PIA servers for now
func (r *reader) GetPortForwarding() (activated bool, err error) {
	s, err := r.envParams.GetEnv("PORT_FORWARDING", libparams.Default("off"))
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
func (r *reader) GetPortForwardingStatusFilepath() (filepath models.Filepath, err error) {
	filepathStr, err := r.envParams.GetPath("PORT_FORWARDING_STATUS_FILE", libparams.Default("/tmp/gluetun/forwarded_port"), libparams.CaseSensitiveValue())
	return models.Filepath(filepathStr), err
}

// GetPIAEncryptionPreset obtains the encryption level for the PIA connection
// from the environment variable PIA_ENCRYPTION, and using ENCRYPTION for
// retro compatibility
func (r *reader) GetPIAEncryptionPreset() (preset string, err error) {
	// Retro-compatibility
	s, err := r.envParams.GetValueIfInside("ENCRYPTION", []string{
		constants.PIAEncryptionPresetNormal,
		constants.PIAEncryptionPresetStrong,
		""})
	if err != nil {
		return "", err
	} else if len(s) != 0 {
		r.logger.Warn("You are using the old environment variable ENCRYPTION, please consider changing it to PIA_ENCRYPTION")
		return s, nil
	}
	return r.envParams.GetValueIfInside(
		"PIA_ENCRYPTION",
		[]string{
			constants.PIAEncryptionPresetNormal,
			constants.PIAEncryptionPresetStrong,
		},
		libparams.Default(constants.PIAEncryptionPresetStrong))
}

// GetPIARegion obtains the region for the PIA server from the
// environment variable REGION
func (r *reader) GetPIARegion() (region string, err error) {
	choices := append(constants.PIAGeoChoices(), "")
	return r.envParams.GetValueIfInside("REGION", choices)
}

// GetPIAOldRegion obtains the region for the PIA server from the
// environment variable REGION
func (r *reader) GetPIAOldRegion() (region string, err error) {
	choices := append(constants.PIAOldGeoChoices(), "")
	return r.envParams.GetValueIfInside("REGION", choices)
}
