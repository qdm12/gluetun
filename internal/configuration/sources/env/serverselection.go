package env

import (
	"errors"
	"fmt"
	"net/netip"
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

	countriesKey, _ := s.getEnvWithRetro("SERVER_COUNTRIES", []string{"COUNTRY"})
	ss.Countries = s.env.CSV(countriesKey)
	if vpnProvider == providers.Cyberghost && len(ss.Countries) == 0 {
		// Retro-compatibility for Cyberghost using the REGION variable
		ss.Countries = s.env.CSV("REGION")
		if len(ss.Countries) > 0 {
			s.onRetroActive("REGION", "SERVER_COUNTRIES")
		}
	}

	regionsKey, _ := s.getEnvWithRetro("SERVER_REGIONS", []string{"REGION"})
	ss.Regions = s.env.CSV(regionsKey)

	citiesKey, _ := s.getEnvWithRetro("SERVER_CITIES", []string{"CITY"})
	ss.Cities = s.env.CSV(citiesKey)

	ss.ISPs = s.env.CSV("ISP")

	hostnamesKey, _ := s.getEnvWithRetro("SERVER_HOSTNAMES", []string{"SERVER_HOSTNAME"})
	ss.Hostnames = s.env.CSV(hostnamesKey)

	serverNamesKey, _ := s.getEnvWithRetro("SERVER_NAMES", []string{"SERVER_NAME"})
	ss.Names = s.env.CSV(serverNamesKey)

	if csv := s.env.Get("SERVER_NUMBER"); csv != nil {
		numbersStrings := strings.Split(*csv, ",")
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
	ss.FreeOnly, err = s.env.BoolPtr("FREE_ONLY")
	if err != nil {
		return ss, err
	}

	// VPNSecure only
	ss.PremiumOnly, err = s.env.BoolPtr("PREMIUM_ONLY")
	if err != nil {
		return ss, err
	}

	// VPNUnlimited only
	ss.MultiHopOnly, err = s.env.BoolPtr("MULTIHOP_ONLY")
	if err != nil {
		return ss, err
	}

	// VPNUnlimited only
	ss.MultiHopOnly, err = s.env.BoolPtr("STREAM_ONLY")
	if err != nil {
		return ss, err
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

func (s *Source) readOpenVPNTargetIP() (ip netip.Addr, err error) {
	envKey, value := s.getEnvWithRetro("VPN_ENDPOINT_IP", []string{"OPENVPN_TARGET_IP"})
	if value == nil {
		return ip, nil
	}

	ip, err = netip.ParseAddr(*value)
	if err != nil {
		return ip, fmt.Errorf("environment variable %s: %w", envKey, err)
	}

	return ip, nil
}

func (s *Source) readOwnedOnly() (ownedOnly *bool, err error) {
	envKey, _ := s.getEnvWithRetro("OWNED_ONLY", []string{"OWNED"})
	return s.env.BoolPtr(envKey)
}
