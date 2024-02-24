package env

import (
	"errors"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gosettings/sources/env"
)

func (s *Source) readServerSelection(vpnProvider, vpnType string) (
	ss settings.ServerSelection, err error) {
	ss.VPN = vpnType

	ss.TargetIP, err = s.env.NetipAddr("VPN_ENDPOINT_IP",
		env.RetroKeys("OPENVPN_TARGET_IP"))
	if err != nil {
		return ss, err
	}

	ss.Countries = s.env.CSV("SERVER_COUNTRIES", env.RetroKeys("COUNTRY"))
	if vpnProvider == providers.Cyberghost && len(ss.Countries) == 0 {
		// Retro-compatibility for Cyberghost using the REGION variable
		ss.Countries = s.env.CSV("REGION")
		if len(ss.Countries) > 0 {
			s.handleDeprecatedKey("REGION", "SERVER_COUNTRIES")
		}
	}

	ss.Regions = s.env.CSV("SERVER_REGIONS", env.RetroKeys("REGION"))
	ss.Cities = s.env.CSV("SERVER_CITIES", env.RetroKeys("CITY"))
	ss.ISPs = s.env.CSV("ISP")
	ss.Hostnames = s.env.CSV("SERVER_HOSTNAMES", env.RetroKeys("SERVER_HOSTNAME"))
	ss.Names = s.env.CSV("SERVER_NAMES", env.RetroKeys("SERVER_NAME"))
	ss.Numbers, err = s.env.CSVUint16("SERVER_NUMBER")
	if err != nil {
		return ss, err
	}

	// Mullvad only
	ss.OwnedOnly, err = s.env.BoolPtr("OWNED_ONLY", env.RetroKeys("OWNED"))
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

	// Surfshark only
	ss.MultiHopOnly, err = s.env.BoolPtr("MULTIHOP_ONLY")
	if err != nil {
		return ss, err
	}

	// VPNUnlimited only
	ss.StreamOnly, err = s.env.BoolPtr("STREAM_ONLY")
	if err != nil {
		return ss, err
	}

	// PIA only
	ss.PortForwardOnly, err = s.env.BoolPtr("PORT_FORWARD_ONLY")
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
