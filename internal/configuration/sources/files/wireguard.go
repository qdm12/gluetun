package files

import (
	"fmt"
	"net/netip"
	"regexp"
	"strings"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"gopkg.in/ini.v1"
)

var (
	regexINISectionNotExist = regexp.MustCompile(`^section ".+" does not exist$`)
	regexINIKeyNotExist     = regexp.MustCompile(`key ".*" not exists$`)
)

func (s *Source) readWireguard() (wireguard settings.Wireguard, err error) {
	fileStringPtr, err := ReadFromFile(s.wireguardConfigPath)
	if err != nil {
		return wireguard, fmt.Errorf("reading file: %w", err)
	}

	if fileStringPtr == nil {
		return wireguard, nil
	}

	rawData := []byte(*fileStringPtr)
	iniFile, err := ini.Load(rawData)
	if err != nil {
		return wireguard, fmt.Errorf("loading ini from reader: %w", err)
	}

	interfaceSection, err := iniFile.GetSection("Interface")
	if err == nil {
		err = parseWireguardInterfaceSection(interfaceSection, &wireguard)
		if err != nil {
			return wireguard, fmt.Errorf("parsing interface section: %w", err)
		}
	} else if !regexINISectionNotExist.MatchString(err.Error()) {
		// can never happen
		return wireguard, fmt.Errorf("getting interface section: %w", err)
	}

	return wireguard, nil
}

func parseWireguardInterfaceSection(interfaceSection *ini.Section,
	wireguard *settings.Wireguard) (err error) {
	wireguard.PrivateKey, err = parseINIWireguardKey(interfaceSection, "PrivateKey")
	if err != nil {
		return err // error is already wrapped correctly
	}

	wireguard.PreSharedKey, err = parseINIWireguardKey(interfaceSection, "PreSharedKey")
	if err != nil {
		return err // error is already wrapped correctly
	}

	wireguard.Addresses, err = parseINIWireguardAddress(interfaceSection)
	if err != nil {
		return err // error is already wrapped correctly
	}

	return nil
}

func parseINIWireguardKey(section *ini.Section, keyName string) (
	key *string, err error) {
	iniKey, err := section.GetKey(keyName)
	if err != nil {
		if regexINIKeyNotExist.MatchString(err.Error()) {
			return nil, nil //nolint:nilnil
		}
		// can never happen
		return nil, fmt.Errorf("getting %s key: %w", keyName, err)
	}

	key = new(string)
	*key = iniKey.String()
	_, err = wgtypes.ParseKey(*key)
	if err != nil {
		return nil, fmt.Errorf("parsing %s: %s: %w", keyName, *key, err)
	}
	return key, nil
}

func parseINIWireguardAddress(section *ini.Section) (
	addresses []netip.Prefix, err error) {
	addressKey, err := section.GetKey("Address")
	if err != nil {
		if regexINIKeyNotExist.MatchString(err.Error()) {
			return nil, nil
		}
		// can never happen
		return nil, fmt.Errorf("getting Address key: %w", err)
	}

	addressStrings := strings.Split(addressKey.String(), ",")
	addresses = make([]netip.Prefix, len(addressStrings))
	for i, addressString := range addressStrings {
		addressString = strings.TrimSpace(addressString)
		if !strings.ContainsRune(addressString, '/') {
			addressString += "/32"
		}
		addresses[i], err = netip.ParsePrefix(addressString)
		if err != nil {
			return nil, fmt.Errorf("parsing address: %w", err)
		}
	}

	return addresses, nil
}
