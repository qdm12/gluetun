package configuration

import (
	"errors"
	"fmt"
	"net"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/golibs/params"
)

var (
	errCustomNotSupported    = errors.New("custom provider is not supported")
	errCustomExtractFromFile = errors.New("cannot extract configuration from file")
)

func (settings *Provider) readCustom(r reader, vpnType string) (err error) {
	settings.Name = constants.Custom

	switch vpnType {
	case constants.OpenVPN:
		return settings.ServerSelection.OpenVPN.readCustom(r)
	case constants.Wireguard:
		return settings.ServerSelection.Wireguard.readCustom(r)
	default:
		return fmt.Errorf("%w: for VPN type %s", errCustomNotSupported, vpnType)
	}
}

func (settings *OpenVPNSelection) readCustom(r reader) (err error) {
	configFile, err := r.env.Get("OPENVPN_CUSTOM_CONFIG", params.CaseSensitiveValue(), params.Compulsory())
	if err != nil {
		return fmt.Errorf("environment variable OPENVPN_CUSTOM_CONFIG: %w", err)
	}
	settings.ConfFile = configFile

	// For display and consistency purposes only,
	// these values are not actually used since the file is re-read
	// before each OpenVPN start.
	_, connection, err := r.ovpnExt.Data(configFile)
	if err != nil {
		return fmt.Errorf("%w: %s", errCustomExtractFromFile, err)
	}
	settings.TCP = connection.Protocol == constants.TCP

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

func (settings *WireguardSelection) readCustom(r reader) (err error) {
	settings.PublicKey, err = r.env.Get("WIREGUARD_PUBLIC_KEY",
		params.CaseSensitiveValue(), params.Compulsory())
	if err != nil {
		return fmt.Errorf("environment variable WIREGUARD_PUBLIC_KEY: %w", err)
	}

	settings.EndpointIP, err = readWireguardEndpointIP(r.env)
	if err != nil {
		return err
	}

	settings.EndpointPort, err = r.env.Port("WIREGUARD_ENDPOINT_PORT", params.Compulsory(),
		params.RetroKeys([]string{"WIREGUARD_PORT"}, r.onRetroActive))
	if err != nil {
		return fmt.Errorf("environment variable WIREGUARD_ENDPOINT_PORT: %w", err)
	}

	return nil
}

// readWireguardEndpointIP reads and parses the server endpoint IP
// address from the environment variable WIREGUARD_ENDPOINT_IP.
func readWireguardEndpointIP(env params.Interface) (endpointIP net.IP, err error) {
	s, err := env.Get("WIREGUARD_ENDPOINT_IP", params.Compulsory())
	if err != nil {
		return nil, fmt.Errorf("environment variable WIREGUARD_ENDPOINT_IP: %w", err)
	}

	endpointIP = net.ParseIP(s)
	if endpointIP == nil {
		return nil, fmt.Errorf("environment variable WIREGUARD_ENDPOINT_IP: %w: %s",
			ErrInvalidIP, s)
	}

	return endpointIP, nil
}
