package params

import (
	"fmt"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	libparams "github.com/qdm12/golibs/params"
)

// GetPortForwarding obtains if port forwarding on the VPN provider server
// side is enabled or not from the environment variable PORT_FORWARDING
// Only valid for older PIA servers for now.
func (r *reader) GetPortForwarding() (activated bool, err error) {
	s, err := r.env.Get("PORT_FORWARDING", libparams.Default("off"))
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
// from the environment variable PORT_FORWARDING_STATUS_FILE.
func (r *reader) GetPortForwardingStatusFilepath() (filepath models.Filepath, err error) {
	filepathStr, err := r.env.Path(
		"PORT_FORWARDING_STATUS_FILE",
		libparams.Default("/tmp/gluetun/forwarded_port"),
		libparams.CaseSensitiveValue())
	return models.Filepath(filepathStr), err
}

// GetPIAEncryptionPreset obtains the encryption level for the PIA connection
// from the environment variable PIA_ENCRYPTION, and using ENCRYPTION for
// retro compatibility.
func (r *reader) GetPIAEncryptionPreset() (preset string, err error) {
	// Retro-compatibility
	s, err := r.env.Inside("ENCRYPTION", []string{
		constants.PIAEncryptionPresetNormal,
		constants.PIAEncryptionPresetStrong})
	if err != nil {
		return "", err
	} else if len(s) != 0 {
		r.logger.Warn("You are using the old environment variable ENCRYPTION, please consider changing it to PIA_ENCRYPTION")
		return s, nil
	}
	return r.env.Inside(
		"PIA_ENCRYPTION",
		[]string{
			constants.PIAEncryptionPresetNormal,
			constants.PIAEncryptionPresetStrong,
		},
		libparams.Default(constants.PIAEncryptionPresetStrong))
}

// GetPIARegions obtains the regions for the PIA servers from the
// environment variable REGION.
func (r *reader) GetPIARegions() (regions []string, err error) {
	return r.env.CSVInside("REGION", constants.PIAGeoChoices())
}

// GetPIAPort obtains the port to reach the PIA server on from the
// environment variable PORT.
func (r *reader) GetPIAPort() (port uint16, err error) {
	n, err := r.env.IntRange("PORT", 0, 65535, libparams.Default("0"))
	return uint16(n), err
}
