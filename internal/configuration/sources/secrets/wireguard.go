package secrets

import (
	"fmt"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/configuration/sources/files"
)

func (s *Source) readWireguard() (settings settings.Wireguard, err error) {
	wireguardConf, err := s.readSecretFileAsStringPtr(
		"WIREGUARD_CONF_SECRETFILE",
		"/run/secrets/wg0.conf",
	)
	if err != nil {
		return settings, fmt.Errorf("reading Wireguard conf secret file: %w", err)
	} else if wireguardConf != nil {
		// Wireguard ini config file takes precedence over individual secrets
		return files.ParseWireguardConf([]byte(*wireguardConf))
	}

	settings.PrivateKey, err = s.readSecretFileAsStringPtr(
		"WIREGUARD_PRIVATE_KEY_SECRETFILE",
		"/run/secrets/wireguard_private_key",
	)
	if err != nil {
		return settings, fmt.Errorf("reading private key file: %w", err)
	}

	settings.PreSharedKey, err = s.readSecretFileAsStringPtr(
		"WIREGUARD_PRESHARED_KEY_SECRETFILE",
		"/run/secrets/wireguard_preshared_key",
	)
	if err != nil {
		return settings, fmt.Errorf("reading preshared key file: %w", err)
	}

	wireguardAddressesCSV, err := s.readSecretFileAsStringPtr(
		"WIREGUARD_ADDRESSES_SECRETFILE",
		"/run/secrets/wireguard_addresses",
	)
	if err != nil {
		return settings, fmt.Errorf("reading addresses file: %w", err)
	} else if wireguardAddressesCSV != nil {
		settings.Addresses, err = parseAddresses(*wireguardAddressesCSV)
		if err != nil {
			return settings, fmt.Errorf("parsing addresses: %w", err)
		}
	}

	return settings, nil
}
