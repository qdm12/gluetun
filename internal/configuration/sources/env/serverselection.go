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

func (r *Reader) readServerSelection(vpnProvider, vpnType string) (
	ss settings.ServerSelection, err error) {
	ss.VPN = vpnType

	ss.TargetIP, err = r.readOpenVPNTargetIP()
	if err != nil {
		return ss, err
	}

	countriesKey, _ := r.getEnvWithRetro("SERVER_COUNTRIES", "COUNTRY")
	ss.Countries = envToCSV(countriesKey)
	if vpnProvider == providers.Cyberghost && len(ss.Countries) == 0 {
		// Retro-compatibility for Cyberghost using the REGION variable
		ss.Countries = envToCSV("REGION")
		if len(ss.Countries) > 0 {
			r.onRetroActive("REGION", "SERVER_COUNTRIES")
		}
	}

	regionsKey, _ := r.getEnvWithRetro("SERVER_REGIONS", "REGION")
	ss.Regions = envToCSV(regionsKey)

	citiesKey, _ := r.getEnvWithRetro("SERVER_CITIES", "CITY")
	ss.Cities = envToCSV(citiesKey)

	ss.ISPs = envToCSV("ISP")

	hostnamesKey, _ := r.getEnvWithRetro("SERVER_HOSTNAMES", "SERVER_HOSTNAME")
	ss.Hostnames = envToCSV(hostnamesKey)

	serverNamesKey, _ := r.getEnvWithRetro("SERVER_NAMES", "SERVER_NAME")
	ss.Names = envToCSV(serverNamesKey)

	if csv := getCleanedEnv("SERVER_NUMBER"); csv != "" {
		numbersStrings := strings.Split(csv, ",")
		numbers := make([]uint16, len(numbersStrings))
		for i, numberString := range numbersStrings {
			number, err := strconv.Atoi(numberString)
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
	ss.OwnedOnly, err = r.readOwnedOnly()
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

	ss.OpenVPN, err = r.readOpenVPNSelection()
	if err != nil {
		return ss, err
	}

	ss.Wireguard, err = r.readWireguardSelection()
	if err != nil {
		return ss, err
	}

	return ss, nil
}

var (
	ErrInvalidIP = errors.New("invalid IP address")
)

func (r *Reader) readOpenVPNTargetIP() (ip net.IP, err error) {
	envKey, s := r.getEnvWithRetro("VPN_ENDPOINT_IP", "OPENVPN_TARGET_IP")
	if s == "" {
		return nil, nil
	}

	ip = net.ParseIP(s)
	if ip == nil {
		return nil, fmt.Errorf("environment variable %s: %w: %s",
			envKey, ErrInvalidIP, s)
	}

	return ip, nil
}

func (r *Reader) readOwnedOnly() (ownedOnly *bool, err error) {
	envKey, _ := r.getEnvWithRetro("OWNED_ONLY", "OWNED")
	ownedOnly, err = envToBoolPtr(envKey)
	if err != nil {
		return nil, fmt.Errorf("environment variable %s: %w", envKey, err)
	}
	return ownedOnly, nil
}
