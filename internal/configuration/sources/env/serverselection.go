package env

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants/providers"
)

var (
	ErrServerNumberNotValid = errors.New("server number is not valid")
)

func (s *Source) readServerSelection(vpnProvider, vpnType string) (
	ss settings.ServerSelection, err error) {
	ss.VPN = vpnType

	ss.TargetIP, err = s.readOpenVPNTargetIP()
	if err != nil {
		return ss, err
	}

	countriesKey, _ := s.getEnvWithRetro("SERVER_COUNTRIES", "COUNTRY")
	ss.Countries = envToCSV(countriesKey)
	if vpnProvider == providers.Cyberghost && len(ss.Countries) == 0 {
		// Retro-compatibility for Cyberghost using the REGION variable
		ss.Countries = envToCSV("REGION")
		if len(ss.Countries) > 0 {
			s.onRetroActive("REGION", "SERVER_COUNTRIES")
		}
	}

	regionsKey, _ := s.getEnvWithRetro("SERVER_REGIONS", "REGION")
	ss.Regions = envToCSV(regionsKey)

	citiesKey, _ := s.getEnvWithRetro("SERVER_CITIES", "CITY")
	ss.Cities = envToCSV(citiesKey)

	ss.ISPs = envToCSV("ISP")

	hostnamesKey, _ := s.getEnvWithRetro("SERVER_HOSTNAMES", "SERVER_HOSTNAME")
	ss.Hostnames = envToCSV(hostnamesKey)

	serverNamesKey, _ := s.getEnvWithRetro("SERVER_NAMES", "SERVER_NAME")
	ss.Names = envToCSV(serverNamesKey)

	if csv := getCleanedEnv("SERVER_NUMBER"); csv != "" {
		numbersStrings := strings.Split(csv, ",")
		numbers := make([]uint16, len(numbersStrings))
		for i, numberString := range numbersStrings {
			const base, bitSize = 10, 16
			number, err := strconv.ParseInt(numberString, base, bitSize)
			if err != nil {
				return ss, fmt.Errorf("%w: %s",
					ErrServerNumberNotValid, numberString)
			} else if number < 0 || number > 65535 {
				return ss, fmt.Errorf("%w: %d must be between 0 and 65535",
					ErrServerNumberNotValid, number)
			}
			numbers[i] = uint16(number)
		}
		ss.Numbers = numbers
	}

	// Mullvad only
	ss.OwnedOnly, err = s.readOwnedOnly()
	if err != nil {
		return ss, err
	}

	// VPNUnlimited and ProtonVPN only
	ss.FreeOnly, err = envToBoolPtr("FREE_ONLY")
	if err != nil {
		return ss, fmt.Errorf("environment variable FREE_ONLY: %w", err)
	}

	// VPNSecure only
	ss.PremiumOnly, err = envToBoolPtr("PREMIUM_ONLY")
	if err != nil {
		return ss, fmt.Errorf("environment variable PREMIUM_ONLY: %w", err)
	}

	// VPNUnlimited only
	ss.MultiHopOnly, err = envToBoolPtr("MULTIHOP_ONLY")
	if err != nil {
		return ss, fmt.Errorf("environment variable MULTIHOP_ONLY: %w", err)
	}

	// VPNUnlimited only
	ss.MultiHopOnly, err = envToBoolPtr("STREAM_ONLY")
	if err != nil {
		return ss, fmt.Errorf("environment variable STREAM_ONLY: %w", err)
	}

	ss.OpenVPN, err = s.readOpenVPNSelection()
	if err != nil {
		return ss, err
	}

	ss.Wireguard, err = s.readWireguardSelection()
	if err != nil {
		return ss, err
	}

	return ss, nil
}

var (
	ErrInvalidIP = errors.New("invalid IP address")
)

func (s *Source) readOpenVPNTargetIP() (ip net.IP, err error) {
	envKey, value := s.getEnvWithRetro("VPN_ENDPOINT_IP", "OPENVPN_TARGET_IP")
	if value == "" {
		return nil, nil
	}

	ip = net.ParseIP(value)
	if ip == nil {
		return nil, fmt.Errorf("environment variable %s: %w: %s",
			envKey, ErrInvalidIP, value)
	}

	return ip, nil
}

func (s *Source) readOwnedOnly() (ownedOnly *bool, err error) {
	envKey, _ := s.getEnvWithRetro("OWNED_ONLY", "OWNED")
	ownedOnly, err = envToBoolPtr(envKey)
	if err != nil {
		return nil, fmt.Errorf("environment variable %s: %w", envKey, err)
	}
	return ownedOnly, nil
}
