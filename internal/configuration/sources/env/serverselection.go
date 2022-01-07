package env

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants"
)

var (
	ErrServerNumberNotValid = errors.New("server number is not valid")
)

func (r *Reader) readServerSelection(vpnProvider, vpnType string) (
	ss settings.ServerSelection, err error) {
	ss.VPN = vpnType

	ss.TargetIP, err = readOpenVPNTargetIP()
	if err != nil {
		return ss, err
	}

	countriesCSV := os.Getenv("COUNTRY")
	if vpnProvider == constants.Cyberghost && countriesCSV == "" {
		// Retro-compatibility
		r.onRetroActive("REGION", "COUNTRY")
		countriesCSV = os.Getenv("REGION")
	}
	if countriesCSV != "" {
		ss.Countries = lowerAndSplit(countriesCSV)
	}

	ss.Regions = envToCSV("REGION")
	ss.Cities = envToCSV("CITY")
	ss.ISPs = envToCSV("ISP")
	ss.Hostnames = envToCSV("SERVER_HOSTNAME")
	ss.Names = envToCSV("SERVER_NAME")

	if csv := os.Getenv("SERVER_NUMBER"); csv != "" {
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
	ss.OwnedOnly, err = envToBoolPtr("OWNED")
	if err != nil {
		return ss, fmt.Errorf("environment variable OWNED: %w", err)
	}

	// VPNUnlimited and ProtonVPN only
	ss.FreeOnly, err = envToBoolPtr("FREE_ONLY")
	if err != nil {
		return ss, fmt.Errorf("environment variable FREE_ONLY: %w", err)
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

func readOpenVPNTargetIP() (ip net.IP, err error) {
	s := os.Getenv("OPENVPN_TARGET_IP")
	if s == "" {
		return nil, nil
	}

	ip = net.ParseIP(s)
	if ip == nil {
		return nil, fmt.Errorf("environment variable OPENVPN_TARGET_IP: %w: %s",
			ErrInvalidIP, s)
	}

	return ip, nil
}
