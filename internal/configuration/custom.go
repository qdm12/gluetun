package configuration

import (
	"errors"
	"fmt"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/golibs/params"
)

var (
	errCustomNotSupported    = errors.New("custom provider is not supported")
	errCustomExtractFromFile = errors.New("cannot extract configuration from file")
)

func (settings *Provider) readCustom(r reader, vpnType string) (err error) {
	settings.Name = constants.Custom

	if vpnType != constants.OpenVPN {
		return fmt.Errorf("%w: for VPN type %s", errCustomNotSupported, vpnType)
	}

	return settings.readCustomOpenVPN(r)
}

func (settings *Provider) readCustomOpenVPN(r reader) (err error) {
	configFile, err := r.env.Get("OPENVPN_CUSTOM_CONFIG", params.CaseSensitiveValue(), params.Compulsory())
	if err != nil {
		return fmt.Errorf("environment variable OPENVPN_CUSTOM_CONFIG: %w", err)
	}
	settings.ServerSelection.OpenVPN.ConfFile = configFile

	// For display and consistency purposes only,
	// these values are not actually used since the file is re-read
	// before each OpenVPN start.
	_, connection, err := r.ovpnExt.Data(configFile)
	if err != nil {
		return fmt.Errorf("%w: %s", errCustomExtractFromFile, err)
	}
	settings.ServerSelection.OpenVPN.TCP = connection.Protocol == constants.TCP

	return nil
}

func (settings *OpenVPN) readCustom(r reader) (err error) {
	settings.ConfFile, err = r.env.Path("OPENVPN_CUSTOM_CONFIG",
		params.Compulsory(), params.CaseSensitiveValue())
	if err != nil {
		return fmt.Errorf("environment variable OPENVPN_CUSTOM_CONFIG: %w", err)
	}

	return nil
}
