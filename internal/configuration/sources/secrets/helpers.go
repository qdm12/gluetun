package secrets

import (
	"fmt"
	"net/netip"
	"strings"

	"github.com/qdm12/gluetun/internal/configuration/sources/files"
	"github.com/qdm12/gluetun/internal/openvpn/extract"
	"github.com/qdm12/gosettings/sources/env"
)

func (s *Source) readSecretFileAsStringPtr(secretPathEnvKey, defaultSecretPath string) (
	stringPtr *string, err error) {
	path := s.env.String(secretPathEnvKey, env.ForceLowercase(false))
	if path == "" {
		path = defaultSecretPath
	}
	return files.ReadFromFile(path)
}

func (s *Source) readPEMSecretFile(secretPathEnvKey, defaultSecretPath string) (
	base64Ptr *string, err error) {
	pemData, err := s.readSecretFileAsStringPtr(secretPathEnvKey, defaultSecretPath)
	if err != nil {
		return nil, fmt.Errorf("reading secret file: %w", err)
	}

	if pemData == nil {
		return nil, nil //nolint:nilnil
	}

	base64Data, err := extract.PEM([]byte(*pemData))
	if err != nil {
		return nil, fmt.Errorf("extracting base64 encoded data from PEM content: %w", err)
	}

	return &base64Data, nil
}

func parseAddresses(addressesCSV string) (addresses []netip.Prefix, err error) {
	if addressesCSV == "" {
		return nil, nil
	}

	addressStrings := strings.Split(addressesCSV, ",")
	addresses = make([]netip.Prefix, len(addressStrings))
	for i, addressString := range addressStrings {
		addressString = strings.TrimSpace(addressString)
		addresses[i], err = netip.ParsePrefix(addressString)
		if err != nil {
			return nil, fmt.Errorf("parsing address %d of %d: %w",
				i+1, len(addressStrings), err)
		}
	}

	return addresses, nil
}
