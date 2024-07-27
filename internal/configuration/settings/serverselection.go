package settings

import (
	"errors"
	"fmt"
	"net/netip"
	"strings"

	"github.com/qdm12/gluetun/internal/configuration/settings/helpers"
	"github.com/qdm12/gluetun/internal/configuration/settings/validation"
	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gluetun/internal/constants/vpn"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gosettings"
	"github.com/qdm12/gosettings/reader"
	"github.com/qdm12/gosettings/validate"
	"github.com/qdm12/gotree"
)

type ServerSelection struct { //nolint:maligned
	// VPN is the VPN type which can be 'openvpn'
	// or 'wireguard'. It cannot be the empty string
	// in the internal state.
	VPN string `json:"vpn"`
	// TargetIP is the server endpoint IP address to use.
	// It will override any IP address from the picked
	// built-in server. It cannot be the empty value in the internal
	// state, and can be set to the unspecified address to indicate
	// there is not target IP address to use.
	TargetIP netip.Addr `json:"target_ip"`
	// Countries is the list of countries to filter VPN servers with.
	Countries []string `json:"countries"`
	// Categories is the list of categories to filter VPN servers with.
	Categories []string `json:"categories"`
	// Regions is the list of regions to filter VPN servers with.
	Regions []string `json:"regions"`
	// Cities is the list of cities to filter VPN servers with.
	Cities []string `json:"cities"`
	// ISPs is the list of ISP names to filter VPN servers with.
	ISPs []string `json:"isps"`
	// Names is the list of server names to filter VPN servers with.
	Names []string `json:"names"`
	// Numbers is the list of server numbers to filter VPN servers with.
	Numbers []uint16 `json:"numbers"`
	// Hostnames is the list of hostnames to filter VPN servers with.
	Hostnames []string `json:"hostnames"`
	// OwnedOnly is true if VPN provider servers that are not owned
	// should be filtered. This is used with Mullvad.
	OwnedOnly *bool `json:"owned_only"`
	// FreeOnly is true if VPN servers that are not free should
	// be filtered. This is used with ProtonVPN and VPN Unlimited.
	FreeOnly *bool `json:"free_only"`
	// PremiumOnly is true if VPN servers that are not premium should
	// be filtered. This is used with VPN Secure.
	// TODO extend to providers using FreeOnly.
	PremiumOnly *bool `json:"premium_only"`
	// StreamOnly is true if VPN servers not for streaming should
	// be filtered. This is used with VPNUnlimited.
	StreamOnly *bool `json:"stream_only"`
	// MultiHopOnly is true if VPN servers that are not multihop
	// should be filtered. This is used with Surfshark.
	MultiHopOnly *bool `json:"multi_hop_only"`
	// PortForwardOnly is true if VPN servers that don't support
	// port forwarding should be filtered. This is used with PIA.
	PortForwardOnly *bool `json:"port_forward_only"`
	// OpenVPN contains settings to select OpenVPN servers
	// and the final connection.
	OpenVPN OpenVPNSelection `json:"openvpn"`
	// Wireguard contains settings to select Wireguard servers
	// and the final connection.
	Wireguard WireguardSelection `json:"wireguard"`
}

var (
	ErrOwnedOnlyNotSupported       = errors.New("owned only filter is not supported")
	ErrFreeOnlyNotSupported        = errors.New("free only filter is not supported")
	ErrPremiumOnlyNotSupported     = errors.New("premium only filter is not supported")
	ErrStreamOnlyNotSupported      = errors.New("stream only filter is not supported")
	ErrMultiHopOnlyNotSupported    = errors.New("multi hop only filter is not supported")
	ErrPortForwardOnlyNotSupported = errors.New("port forwarding only filter is not supported")
	ErrFreePremiumBothSet          = errors.New("free only and premium only filters are both set")
)

func (ss *ServerSelection) validate(vpnServiceProvider string,
	storage Storage) (err error) {
	switch ss.VPN {
	case vpn.OpenVPN, vpn.Wireguard:
	default:
		return fmt.Errorf("%w: %s", ErrVPNTypeNotValid, ss.VPN)
	}

	filterChoices, err := getLocationFilterChoices(vpnServiceProvider, ss, storage)
	if err != nil {
		return err // already wrapped error
	}

	// Retro-compatibility
	switch vpnServiceProvider {
	case providers.Nordvpn:
		*ss = nordvpnRetroRegion(*ss, filterChoices.Regions, filterChoices.Countries)
	case providers.Surfshark:
		*ss = surfsharkRetroRegion(*ss)
	}

	err = validateServerFilters(*ss, filterChoices, vpnServiceProvider)
	if err != nil {
		return fmt.Errorf("for VPN service provider %s: %w", vpnServiceProvider, err)
	}

	if *ss.OwnedOnly &&
		vpnServiceProvider != providers.Mullvad {
		return fmt.Errorf("%w: for VPN service provider %s",
			ErrOwnedOnlyNotSupported, vpnServiceProvider)
	}

	if *ss.FreeOnly &&
		!helpers.IsOneOf(vpnServiceProvider,
			providers.Protonvpn,
			providers.VPNUnlimited,
		) {
		return fmt.Errorf("%w: for VPN service provider %s",
			ErrFreeOnlyNotSupported, vpnServiceProvider)
	}

	if *ss.PremiumOnly &&
		!helpers.IsOneOf(vpnServiceProvider,
			providers.VPNSecure,
		) {
		return fmt.Errorf("%w: for VPN service provider %s",
			ErrPremiumOnlyNotSupported, vpnServiceProvider)
	}

	if *ss.FreeOnly && *ss.PremiumOnly {
		return fmt.Errorf("%w", ErrFreePremiumBothSet)
	}

	if *ss.StreamOnly &&
		!helpers.IsOneOf(vpnServiceProvider,
			providers.Protonvpn,
			providers.VPNUnlimited,
		) {
		return fmt.Errorf("%w: for VPN service provider %s",
			ErrStreamOnlyNotSupported, vpnServiceProvider)
	}

	if *ss.MultiHopOnly &&
		vpnServiceProvider != providers.Surfshark {
		return fmt.Errorf("%w: for VPN service provider %s",
			ErrMultiHopOnlyNotSupported, vpnServiceProvider)
	}

	if *ss.PortForwardOnly &&
		vpnServiceProvider != providers.PrivateInternetAccess {
		// ProtonVPN also supports port forwarding, but on all their servers, so these
		// don't have the port forwarding boolean field. As a consequence, we only allow
		// the use of PortForwardOnly for Private Internet Access.
		return fmt.Errorf("%w: for VPN service provider %s",
			ErrPortForwardOnlyNotSupported, vpnServiceProvider)
	}

	if ss.VPN == vpn.OpenVPN {
		err = ss.OpenVPN.validate(vpnServiceProvider)
		if err != nil {
			return fmt.Errorf("OpenVPN server selection settings: %w", err)
		}
	} else {
		err = ss.Wireguard.validate(vpnServiceProvider)
		if err != nil {
			return fmt.Errorf("Wireguard server selection settings: %w", err)
		}
	}

	return nil
}

func getLocationFilterChoices(vpnServiceProvider string,
	ss *ServerSelection, storage Storage) (filterChoices models.FilterChoices,
	err error) {
	filterChoices = storage.GetFilterChoices(vpnServiceProvider)

	if vpnServiceProvider == providers.Surfshark {
		// // Retro compatibility
		// TODO v4 remove
		newAndRetroRegions := append(filterChoices.Regions, validation.SurfsharkRetroLocChoices()...) //nolint:gocritic
		err := validate.AreAllOneOfCaseInsensitive(ss.Regions, newAndRetroRegions)
		if err != nil {
			// Only return error comparing with newer regions, we don't want to confuse the user
			// with the retro regions in the error message.
			err = validate.AreAllOneOfCaseInsensitive(ss.Regions, filterChoices.Regions)
			return models.FilterChoices{}, fmt.Errorf("%w: %w", ErrRegionNotValid, err)
		}
	}

	return filterChoices, nil
}

// validateServerFilters validates filters against the choices given as arguments.
// Set an argument to nil to pass the check for a particular filter.
func validateServerFilters(settings ServerSelection, filterChoices models.FilterChoices,
	vpnServiceProvider string) (err error) {
	err = validate.AreAllOneOfCaseInsensitive(settings.Countries, filterChoices.Countries)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCountryNotValid, err)
	}

	err = validate.AreAllOneOfCaseInsensitive(settings.Regions, filterChoices.Regions)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrRegionNotValid, err)
	}

	err = validate.AreAllOneOfCaseInsensitive(settings.Cities, filterChoices.Cities)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCityNotValid, err)
	}

	err = validate.AreAllOneOfCaseInsensitive(settings.ISPs, filterChoices.ISPs)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrISPNotValid, err)
	}

	err = validate.AreAllOneOfCaseInsensitive(settings.Hostnames, filterChoices.Hostnames)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrHostnameNotValid, err)
	}

	if vpnServiceProvider == providers.Custom && len(settings.Names) == 1 {
		// Allow a single name to be specified for the custom provider in case
		// the user wants to use VPN server side port forwarding with PIA
		// which requires a server name for TLS verification.
		filterChoices.Names = settings.Names
	}
	err = validate.AreAllOneOfCaseInsensitive(settings.Names, filterChoices.Names)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrNameNotValid, err)
	}

	err = validate.AreAllOneOfCaseInsensitive(settings.Categories, filterChoices.Categories)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCategoryNotValid, err)
	}

	return nil
}

func (ss *ServerSelection) copy() (copied ServerSelection) {
	return ServerSelection{
		VPN:             ss.VPN,
		TargetIP:        ss.TargetIP,
		Countries:       gosettings.CopySlice(ss.Countries),
		Categories:      gosettings.CopySlice(ss.Categories),
		Regions:         gosettings.CopySlice(ss.Regions),
		Cities:          gosettings.CopySlice(ss.Cities),
		ISPs:            gosettings.CopySlice(ss.ISPs),
		Hostnames:       gosettings.CopySlice(ss.Hostnames),
		Names:           gosettings.CopySlice(ss.Names),
		Numbers:         gosettings.CopySlice(ss.Numbers),
		OwnedOnly:       gosettings.CopyPointer(ss.OwnedOnly),
		FreeOnly:        gosettings.CopyPointer(ss.FreeOnly),
		PremiumOnly:     gosettings.CopyPointer(ss.PremiumOnly),
		StreamOnly:      gosettings.CopyPointer(ss.StreamOnly),
		PortForwardOnly: gosettings.CopyPointer(ss.PortForwardOnly),
		MultiHopOnly:    gosettings.CopyPointer(ss.MultiHopOnly),
		OpenVPN:         ss.OpenVPN.copy(),
		Wireguard:       ss.Wireguard.copy(),
	}
}

func (ss *ServerSelection) overrideWith(other ServerSelection) {
	ss.VPN = gosettings.OverrideWithComparable(ss.VPN, other.VPN)
	ss.TargetIP = gosettings.OverrideWithValidator(ss.TargetIP, other.TargetIP)
	ss.Countries = gosettings.OverrideWithSlice(ss.Countries, other.Countries)
	ss.Categories = gosettings.OverrideWithSlice(ss.Categories, other.Categories)
	ss.Regions = gosettings.OverrideWithSlice(ss.Regions, other.Regions)
	ss.Cities = gosettings.OverrideWithSlice(ss.Cities, other.Cities)
	ss.ISPs = gosettings.OverrideWithSlice(ss.ISPs, other.ISPs)
	ss.Hostnames = gosettings.OverrideWithSlice(ss.Hostnames, other.Hostnames)
	ss.Names = gosettings.OverrideWithSlice(ss.Names, other.Names)
	ss.Numbers = gosettings.OverrideWithSlice(ss.Numbers, other.Numbers)
	ss.OwnedOnly = gosettings.OverrideWithPointer(ss.OwnedOnly, other.OwnedOnly)
	ss.FreeOnly = gosettings.OverrideWithPointer(ss.FreeOnly, other.FreeOnly)
	ss.PremiumOnly = gosettings.OverrideWithPointer(ss.PremiumOnly, other.PremiumOnly)
	ss.StreamOnly = gosettings.OverrideWithPointer(ss.StreamOnly, other.StreamOnly)
	ss.MultiHopOnly = gosettings.OverrideWithPointer(ss.MultiHopOnly, other.MultiHopOnly)
	ss.PortForwardOnly = gosettings.OverrideWithPointer(ss.PortForwardOnly, other.PortForwardOnly)
	ss.OpenVPN.overrideWith(other.OpenVPN)
	ss.Wireguard.overrideWith(other.Wireguard)
}

func (ss *ServerSelection) setDefaults(vpnProvider string) {
	ss.VPN = gosettings.DefaultComparable(ss.VPN, vpn.OpenVPN)
	ss.TargetIP = gosettings.DefaultValidator(ss.TargetIP, netip.IPv4Unspecified())
	ss.OwnedOnly = gosettings.DefaultPointer(ss.OwnedOnly, false)
	ss.FreeOnly = gosettings.DefaultPointer(ss.FreeOnly, false)
	ss.PremiumOnly = gosettings.DefaultPointer(ss.PremiumOnly, false)
	ss.StreamOnly = gosettings.DefaultPointer(ss.StreamOnly, false)
	ss.MultiHopOnly = gosettings.DefaultPointer(ss.MultiHopOnly, false)
	ss.PortForwardOnly = gosettings.DefaultPointer(ss.PortForwardOnly, false)
	ss.OpenVPN.setDefaults(vpnProvider)
	ss.Wireguard.setDefaults()
}

func (ss ServerSelection) String() string {
	return ss.toLinesNode().String()
}

func (ss ServerSelection) toLinesNode() (node *gotree.Node) {
	node = gotree.New("Server selection settings:")
	node.Appendf("VPN type: %s", ss.VPN)
	if !ss.TargetIP.IsUnspecified() {
		node.Appendf("Target IP address: %s", ss.TargetIP)
	}

	if len(ss.Countries) > 0 {
		node.Appendf("Countries: %s", strings.Join(ss.Countries, ", "))
	}

	if len(ss.Categories) > 0 {
		node.Appendf("Categories: %s", strings.Join(ss.Categories, ", "))
	}

	if len(ss.Regions) > 0 {
		node.Appendf("Regions: %s", strings.Join(ss.Regions, ", "))
	}

	if len(ss.Cities) > 0 {
		node.Appendf("Cities: %s", strings.Join(ss.Cities, ", "))
	}

	if len(ss.ISPs) > 0 {
		node.Appendf("ISPs: %s", strings.Join(ss.ISPs, ", "))
	}

	if len(ss.Names) > 0 {
		node.Appendf("Server names: %s", strings.Join(ss.Names, ", "))
	}

	if len(ss.Numbers) > 0 {
		numbersNode := node.Appendf("Server numbers:")
		for _, number := range ss.Numbers {
			numbersNode.Appendf("%d", number)
		}
	}

	if len(ss.Hostnames) > 0 {
		node.Appendf("Hostnames: %s", strings.Join(ss.Hostnames, ", "))
	}

	if *ss.OwnedOnly {
		node.Appendf("Owned only servers: yes")
	}

	if *ss.FreeOnly {
		node.Appendf("Free only servers: yes")
	}

	if *ss.PremiumOnly {
		node.Appendf("Premium only servers: yes")
	}

	if *ss.StreamOnly {
		node.Appendf("Stream only servers: yes")
	}

	if *ss.MultiHopOnly {
		node.Appendf("Multi-hop only servers: yes")
	}

	if ss.VPN == vpn.OpenVPN {
		node.AppendNode(ss.OpenVPN.toLinesNode())
	} else {
		node.AppendNode(ss.Wireguard.toLinesNode())
	}

	return node
}

// WithDefaults is a shorthand using setDefaults.
// It's used in unit tests in other packages.
func (ss ServerSelection) WithDefaults(provider string) ServerSelection {
	ss.setDefaults(provider)
	return ss
}

func (ss *ServerSelection) read(r *reader.Reader,
	vpnProvider, vpnType string) (err error) {
	ss.VPN = vpnType

	ss.TargetIP, err = r.NetipAddr("OPENVPN_ENDPOINT_IP",
		reader.RetroKeys("OPENVPN_TARGET_IP", "VPN_ENDPOINT_IP"))
	if err != nil {
		return err
	}

	countriesRetroKeys := []string{"COUNTRY"}
	if vpnProvider == providers.Cyberghost {
		countriesRetroKeys = append(countriesRetroKeys, "REGION")
	}
	ss.Countries = r.CSV("SERVER_COUNTRIES", reader.RetroKeys(countriesRetroKeys...))

	ss.Regions = r.CSV("SERVER_REGIONS", reader.RetroKeys("REGION"))
	ss.Cities = r.CSV("SERVER_CITIES", reader.RetroKeys("CITY"))
	ss.ISPs = r.CSV("ISP")
	ss.Hostnames = r.CSV("SERVER_HOSTNAMES", reader.RetroKeys("SERVER_HOSTNAME"))
	ss.Names = r.CSV("SERVER_NAMES", reader.RetroKeys("SERVER_NAME"))
	ss.Numbers, err = r.CSVUint16("SERVER_NUMBER")
	ss.Categories = r.CSV("SERVER_CATEGORIES")
	if err != nil {
		return err
	}

	// Mullvad only
	ss.OwnedOnly, err = r.BoolPtr("OWNED_ONLY", reader.RetroKeys("OWNED"))
	if err != nil {
		return err
	}

	// VPNUnlimited and ProtonVPN only
	ss.FreeOnly, err = r.BoolPtr("FREE_ONLY")
	if err != nil {
		return err
	}

	// VPNSecure only
	ss.PremiumOnly, err = r.BoolPtr("PREMIUM_ONLY")
	if err != nil {
		return err
	}

	// Surfshark only
	ss.MultiHopOnly, err = r.BoolPtr("MULTIHOP_ONLY")
	if err != nil {
		return err
	}

	// VPNUnlimited only
	ss.StreamOnly, err = r.BoolPtr("STREAM_ONLY")
	if err != nil {
		return err
	}

	// PIA only
	ss.PortForwardOnly, err = r.BoolPtr("PORT_FORWARD_ONLY")
	if err != nil {
		return err
	}

	err = ss.OpenVPN.read(r)
	if err != nil {
		return err
	}

	err = ss.Wireguard.read(r)
	if err != nil {
		return err
	}

	return nil
}
