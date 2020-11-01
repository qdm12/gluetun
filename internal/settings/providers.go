package settings

import (
	"fmt"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/params"
)

// GetPIASettings obtains PIA settings from environment variables using the params package.
func GetPIASettings(paramsReader params.Reader) (settings models.ProviderSettings, err error) {
	return getPIASettings(paramsReader, constants.PrivateInternetAccess)
}

// GetPIAOldSettings obtains PIA settings for the older PIA servers (pre summer 2020)
// from environment variables using the params package.
func GetPIAOldSettings(paramsReader params.Reader) (settings models.ProviderSettings, err error) {
	return getPIASettings(paramsReader, constants.PrivateInternetAccessOld)
}

func getPIASettings(paramsReader params.Reader, name models.VPNProvider) (settings models.ProviderSettings, err error) {
	settings.Name = name
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
	settings.ServerSelection.Regions, err = paramsReader.GetPIARegions()
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
	settings.ServerSelection.Countries, err = paramsReader.GetMullvadCountries()
	if err != nil {
		return settings, err
	}
	settings.ServerSelection.Cities, err = paramsReader.GetMullvadCities()
	if err != nil {
		return settings, err
	}
	settings.ServerSelection.ISPs, err = paramsReader.GetMullvadISPs()
	if err != nil {
		return settings, err
	}
	settings.ServerSelection.CustomPort, err = paramsReader.GetMullvadPort()
	if err != nil {
		return settings, err
	}
	if settings.ServerSelection.Protocol == constants.TCP {
		switch settings.ServerSelection.CustomPort {
		case 0, 80, 443, 1401: //nolint:gomnd
		default:
			return settings, fmt.Errorf("port %d is not valid for TCP protocol", settings.ServerSelection.CustomPort)
		}
	} else {
		switch settings.ServerSelection.CustomPort {
		case 0, 53, 1194, 1195, 1196, 1197, 1300, 1301, 1302, 1303, 1400: //nolint:gomnd
		default:
			return settings, fmt.Errorf("port %d is not valid for UDP protocol", settings.ServerSelection.CustomPort)
		}
	}
	settings.ServerSelection.Owned, err = paramsReader.GetMullvadOwned()
	if err != nil {
		return settings, err
	}
	settings.ExtraConfigOptions.OpenVPNIPv6, err = paramsReader.GetOpenVPNIPv6()
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
	settings.ServerSelection.Regions, err = paramsReader.GetWindscribeRegions()
	if err != nil {
		return settings, err
	}
	settings.ServerSelection.Cities, err = paramsReader.GetWindscribeCities()
	if err != nil {
		return settings, err
	}
	settings.ServerSelection.Hostnames, err = paramsReader.GetWindscribeHostnames()
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
	settings.ServerSelection.Regions, err = paramsReader.GetSurfsharkRegions()
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
	settings.ServerSelection.Regions, err = paramsReader.GetCyberghostRegions()
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
	settings.ServerSelection.Regions, err = paramsReader.GetVyprvpnRegions()
	if err != nil {
		return settings, err
	}
	return settings, nil
}

// GetNordvpnSettings obtains NordVPN settings from environment variables using the params package.
func GetNordvpnSettings(paramsReader params.Reader) (settings models.ProviderSettings, err error) {
	settings.Name = constants.Nordvpn
	settings.ServerSelection.Protocol, err = paramsReader.GetNetworkProtocol()
	if err != nil {
		return settings, err
	}
	settings.ServerSelection.TargetIP, err = paramsReader.GetTargetIP()
	if err != nil {
		return settings, err
	}
	settings.ServerSelection.Regions, err = paramsReader.GetNordvpnRegions()
	if err != nil {
		return settings, err
	}
	settings.ServerSelection.Numbers, err = paramsReader.GetNordvpnNumbers()
	if err != nil {
		return settings, err
	}
	return settings, nil
}

// GetPurevpnSettings obtains Purevpn settings from environment variables using the params package.
func GetPurevpnSettings(paramsReader params.Reader) (settings models.ProviderSettings, err error) {
	settings.Name = constants.Mullvad
	settings.ServerSelection.Protocol, err = paramsReader.GetNetworkProtocol()
	if err != nil {
		return settings, err
	}
	settings.ServerSelection.TargetIP, err = paramsReader.GetTargetIP()
	if err != nil {
		return settings, err
	}
	settings.ServerSelection.Regions, err = paramsReader.GetPurevpnRegions()
	if err != nil {
		return settings, err
	}
	settings.ServerSelection.Countries, err = paramsReader.GetPurevpnCountries()
	if err != nil {
		return settings, err
	}
	settings.ServerSelection.Cities, err = paramsReader.GetPurevpnCities()
	if err != nil {
		return settings, err
	}
	return settings, nil
}
