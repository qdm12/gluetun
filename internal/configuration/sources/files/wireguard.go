package files

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/ini.v1"
)

func (s *Source) lazyLoadWireguardConf() WireguardConfig {
	if s.cached.wireguardLoaded {
		return s.cached.wireguardConf
	}

	s.cached.wireguardLoaded = true
	var err error
	s.cached.wireguardConf, err = ParseWireguardConf(filepath.Join(s.rootDirectory, "wireguard", "wg0.conf"))
	if err != nil {
		s.warner.Warnf("skipping Wireguard config: %s", err)
	}
	return s.cached.wireguardConf
}

type AmneziaWgConfig struct {
	Jc   *string
	Jmin *string
	Jmax *string
	S1   *string
	S2   *string
	S3   *string
	S4   *string
	H1   *string
	H2   *string
	H3   *string
	H4   *string
	I1   *string
	I2   *string
	I3   *string
	I4   *string
	I5   *string
}

func NewAmneziaWgConfigFromInerfaceIniSection(interfaceSection *ini.Section) *AmneziaWgConfig {
	return &AmneziaWgConfig{
		Jc:   getINIKeyFromSection(interfaceSection, "Jc"),
		Jmin: getINIKeyFromSection(interfaceSection, "Jmin"),
		Jmax: getINIKeyFromSection(interfaceSection, "Jmax"),
		S1:   getINIKeyFromSection(interfaceSection, "S1"),
		S2:   getINIKeyFromSection(interfaceSection, "S2"),
		S3:   getINIKeyFromSection(interfaceSection, "S3"),
		S4:   getINIKeyFromSection(interfaceSection, "S4"),
		H1:   getINIKeyFromSection(interfaceSection, "H1"),
		H2:   getINIKeyFromSection(interfaceSection, "H2"),
		H3:   getINIKeyFromSection(interfaceSection, "H3"),
		H4:   getINIKeyFromSection(interfaceSection, "H4"),
		I1:   getINIKeyFromSection(interfaceSection, "I1"),
		I2:   getINIKeyFromSection(interfaceSection, "I2"),
		I3:   getINIKeyFromSection(interfaceSection, "I3"),
		I4:   getINIKeyFromSection(interfaceSection, "I4"),
		I5:   getINIKeyFromSection(interfaceSection, "I5"),
	}
}

type WireguardConfig struct {
	PrivateKey    *string
	PreSharedKey  *string
	Addresses     *string
	PublicKey     *string
	EndpointIP    *string
	EndpointPort  *string
	AmneziaParams AmneziaWgConfig
}

var regexINISectionNotExist = regexp.MustCompile(`^section ".+" does not exist$`)

func ParseWireguardConf(path string) (config WireguardConfig, err error) {
	iniFile, err := ini.InsensitiveLoad(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return WireguardConfig{}, nil
		}
		return WireguardConfig{}, fmt.Errorf("loading ini from reader: %w", err)
	}

	interfaceSection, err := iniFile.GetSection("Interface")
	if err == nil {
		config.PrivateKey, config.Addresses = parseWireguardInterfaceSection(interfaceSection)
		config.AmneziaParams = *NewAmneziaWgConfigFromInerfaceIniSection(interfaceSection)
	} else if !regexINISectionNotExist.MatchString(err.Error()) {
		// can never happen
		return WireguardConfig{}, fmt.Errorf("getting interface section: %w", err)
	}

	peerSection, err := iniFile.GetSection("Peer")
	if err == nil {
		config.PreSharedKey, config.PublicKey, config.EndpointIP,
			config.EndpointPort = parseWireguardPeerSection(peerSection)
	} else if !regexINISectionNotExist.MatchString(err.Error()) {
		// can never happen
		return WireguardConfig{}, fmt.Errorf("getting peer section: %w", err)
	}

	return config, nil
}

func parseWireguardInterfaceSection(interfaceSection *ini.Section) (privateKey, addresses *string) {
	privateKey = getINIKeyFromSection(interfaceSection, "PrivateKey")
	addresses = getINIKeyFromSection(interfaceSection, "Address")
	return privateKey, addresses
}

var ErrEndpointHostNotIP = errors.New("endpoint host is not an IP")

func parseWireguardPeerSection(peerSection *ini.Section) (
	preSharedKey, publicKey, endpointIP, endpointPort *string,
) {
	preSharedKey = getINIKeyFromSection(peerSection, "PresharedKey")
	publicKey = getINIKeyFromSection(peerSection, "PublicKey")
	endpoint := getINIKeyFromSection(peerSection, "Endpoint")
	if endpoint != nil {
		parts := strings.Split(*endpoint, ":")
		endpointIP = &parts[0]
		const partsWithPort = 2
		if len(parts) >= partsWithPort {
			endpointPort = new(string)
			*endpointPort = strings.Join(parts[1:], ":")
		}
	}

	return preSharedKey, publicKey, endpointIP, endpointPort
}

var regexINIKeyNotExist = regexp.MustCompile(`key ".*" not exists$`)

func getINIKeyFromSection(section *ini.Section, key string) (value *string) {
	iniKey, err := section.GetKey(key)
	if err != nil {
		if regexINIKeyNotExist.MatchString(err.Error()) {
			return nil
		}
		// can never happen
		panic(fmt.Sprintf("getting key %q: %s", key, err))
	}
	value = new(string)
	*value = iniKey.String()
	return value
}
