package settings

import (
	"github.com/qdm12/private-internet-access-docker/internal/constants"
	"github.com/qdm12/private-internet-access-docker/internal/models"
	"github.com/qdm12/private-internet-access-docker/internal/params"
)

// GetPIASettings obtains PIA settings from environment variables using the params package.
func GetPIASettings(paramsReader params.Reader) (settings models.ProviderSettings, err error) {
	settings.Name = constants.PrivateInternetAccess
	settings.ServerSelection.Protocol, err = paramsReader.GetNetworkProtocol()
	if err != nil {
		return settings, err
	}
	settings.ServerSelection.TargetIP, err = paramsReader.GetTargetIP()
	if err != nil {
		return settings, err
	}
	encryptionPreset, err := paramsReader.GetPIAEncryptionPreset()
	if err != nil {
		return settings, err
	}
	settings.ServerSelection.EncryptionPreset = encryptionPreset
	settings.ExtraConfigOptions.EncryptionPreset = encryptionPreset
	settings.ServerSelection.Region, err = paramsReader.GetPIARegion()
	if err != nil {
		return settings, err
	}
	settings.PortForwarding.Enabled, err = paramsReader.GetPortForwarding()
	if err != nil {
		return settings, err
	}
	if settings.PortForwarding.Enabled {
		settings.PortForwarding.Filepath, err = paramsReader.GetPortForwardingStatusFilepath()
		if err != nil {
			return settings, err
		}
	}
	return settings, nil
}

// GetMullvadSettings obtains Mullvad settings from environment variables using the params package.
func GetMullvadSettings(paramsReader params.Reader) (settings models.ProviderSettings, err error) {
	settings.Name = constants.Mullvad
	settings.ServerSelection.Protocol, err = paramsReader.GetNetworkProtocol()
	if err != nil {
		return settings, err
	}
	settings.ServerSelection.TargetIP, err = paramsReader.GetTargetIP()
	if err != nil {
		return settings, err
	}
	settings.ServerSelection.Country, err = paramsReader.GetMullvadCountry()
	if err != nil {
		return settings, err
	}
	settings.ServerSelection.City, err = paramsReader.GetMullvadCity()
	if err != nil {
		return settings, err
	}
	settings.ServerSelection.ISP, err = paramsReader.GetMullvadISP()
	if err != nil {
		return settings, err
	}
	settings.ServerSelection.CustomPort, err = paramsReader.GetMullvadPort()
	if err != nil {
		return settings, err
	}
	return settings, nil
}

// GetWindscribeSettings obtains Windscribe settings from environment variables using the params package.
func GetWindscribeSettings(paramsReader params.Reader) (settings models.ProviderSettings, err error) {
	settings.Name = constants.Windscribe
	settings.ServerSelection.Protocol, err = paramsReader.GetNetworkProtocol()
	if err != nil {
		return settings, err
	}
	settings.ServerSelection.TargetIP, err = paramsReader.GetTargetIP()
	if err != nil {
		return settings, err
	}
	settings.ServerSelection.Region, err = paramsReader.GetWindscribeRegion()
	if err != nil {
		return settings, err
	}
	settings.ServerSelection.CustomPort, err = paramsReader.GetWindscribePort(settings.ServerSelection.Protocol)
	if err != nil {
		return settings, err
	}
	return settings, nil
}

// GetSurfsharkSettings obtains Surfshark settings from environment variables using the params package.
func GetSurfsharkSettings(paramsReader params.Reader) (settings models.ProviderSettings, err error) {
	settings.Name = constants.Surfshark
	settings.ServerSelection.Protocol, err = paramsReader.GetNetworkProtocol()
	if err != nil {
		return settings, err
	}
	settings.ServerSelection.TargetIP, err = paramsReader.GetTargetIP()
	if err != nil {
		return settings, err
	}
	settings.ServerSelection.Region, err = paramsReader.GetSurfsharkRegion()
	if err != nil {
		return settings, err
	}
	return settings, nil
}

// GetCyberghostSettings obtains Cyberghost settings from environment variables using the params package.
func GetCyberghostSettings(paramsReader params.Reader) (settings models.ProviderSettings, err error) {
	settings.Name = constants.Cyberghost
	settings.ServerSelection.Protocol, err = paramsReader.GetNetworkProtocol()
	if err != nil {
		return settings, err
	}
	settings.ServerSelection.TargetIP, err = paramsReader.GetTargetIP()
	if err != nil {
		return settings, err
	}
	settings.ExtraConfigOptions.ClientKey, err = paramsReader.GetCyberghostClientKey()
	if err != nil {
		return settings, err
	}
	settings.ServerSelection.Group, err = paramsReader.GetCyberghostGroup()
	if err != nil {
		return settings, err
	}
	settings.ServerSelection.Region, err = paramsReader.GetCyberghostRegion()
	if err != nil {
		return settings, err
	}
	return settings, nil
}

// GetVyprvpnSettings obtains Vyprvpn settings from environment variables using the params package.
func GetVyprvpnSettings(paramsReader params.Reader) (settings models.ProviderSettings, err error) {
	settings.Name = constants.Vyprvpn
	settings.ServerSelection.Protocol, err = paramsReader.GetNetworkProtocol()
	if err != nil {
		return settings, err
	}
	settings.ServerSelection.TargetIP, err = paramsReader.GetTargetIP()
	if err != nil {
		return settings, err
	}
	settings.ServerSelection.Region, err = paramsReader.GetVyprvpnRegion()
	if err != nil {
		return settings, err
	}
	return settings, nil
}
