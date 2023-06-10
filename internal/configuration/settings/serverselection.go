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
	// Counties is the list of countries to filter VPN servers with.
	Countries []string `json:"countries"`
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

	// OpenVPN contains settings to select OpenVPN servers
	// and the final connection.
	OpenVPN OpenVPNSelection `json:"openvpn"`
	// Wireguard contains settings to select Wireguard servers
	// and the final connection.
	Wireguard WireguardSelection `json:"wireguard"`
}

var (
	ErrOwnedOnlyNotSupported    = errors.New("owned only filter is not supported")
	ErrFreeOnlyNotSupported     = errors.New("free only filter is not supported")
	ErrPremiumOnlyNotSupported  = errors.New("premium only filter is not supported")
	ErrStreamOnlyNotSupported   = errors.New("stream only filter is not supported")
	ErrMultiHopOnlyNotSupported = errors.New("multi hop only filter is not supported")
	ErrFreePremiumBothSet       = errors.New("free only and premium only filters are both set")
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

	err = validateServerFilters(*ss, filterChoices)
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
		filterChoices.Regions = append(filterChoices.Regions, validation.SurfsharkRetroLocChoices()...)
		err := validate.AreAllOneOfCaseInsensitive(ss.Regions, filterChoices.Regions)
		if err != nil {
			return models.FilterChoices{}, fmt.Errorf("%w: %w", ErrRegionNotValid, err)
		}
	}

	return filterChoices, nil
}

// validateServerFilters validates filters against the choices given as arguments.
// Set an argument to nil to pass the check for a particular filter.
func validateServerFilters(settings ServerSelection, filterChoices models.FilterChoices) (err error) {
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

	err = validate.AreAllOneOfCaseInsensitive(settings.Names, filterChoices.Names)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrNameNotValid, err)
	}

	return nil
}

func (ss *ServerSelection) copy() (copied ServerSelection) {
	return ServerSelection{
		VPN:          ss.VPN,
		TargetIP:     ss.TargetIP,
		Countries:    gosettings.CopySlice(ss.Countries),
		Regions:      gosettings.CopySlice(ss.Regions),
		Cities:       gosettings.CopySlice(ss.Cities),
		ISPs:         gosettings.CopySlice(ss.ISPs),
		Hostnames:    gosettings.CopySlice(ss.Hostnames),
		Names:        gosettings.CopySlice(ss.Names),
		Numbers:      gosettings.CopySlice(ss.Numbers),
		OwnedOnly:    gosettings.CopyPointer(ss.OwnedOnly),
		FreeOnly:     gosettings.CopyPointer(ss.FreeOnly),
		PremiumOnly:  gosettings.CopyPointer(ss.PremiumOnly),
		StreamOnly:   gosettings.CopyPointer(ss.StreamOnly),
		MultiHopOnly: gosettings.CopyPointer(ss.MultiHopOnly),
		OpenVPN:      ss.OpenVPN.copy(),
		Wireguard:    ss.Wireguard.copy(),
	}
}

func (ss *ServerSelection) mergeWith(other ServerSelection) {
	ss.VPN = gosettings.MergeWithString(ss.VPN, other.VPN)
	ss.TargetIP = gosettings.MergeWithValidator(ss.TargetIP, other.TargetIP)
	ss.Countries = gosettings.MergeWithSlice(ss.Countries, other.Countries)
	ss.Regions = gosettings.MergeWithSlice(ss.Regions, other.Regions)
	ss.Cities = gosettings.MergeWithSlice(ss.Cities, other.Cities)
	ss.ISPs = gosettings.MergeWithSlice(ss.ISPs, other.ISPs)
	ss.Hostnames = gosettings.MergeWithSlice(ss.Hostnames, other.Hostnames)
	ss.Names = gosettings.MergeWithSlice(ss.Names, other.Names)
	ss.Numbers = gosettings.MergeWithSlice(ss.Numbers, other.Numbers)
	ss.OwnedOnly = gosettings.MergeWithPointer(ss.OwnedOnly, other.OwnedOnly)
	ss.FreeOnly = gosettings.MergeWithPointer(ss.FreeOnly, other.FreeOnly)
	ss.PremiumOnly = gosettings.MergeWithPointer(ss.PremiumOnly, other.PremiumOnly)
	ss.StreamOnly = gosettings.MergeWithPointer(ss.StreamOnly, other.StreamOnly)
	ss.MultiHopOnly = gosettings.MergeWithPointer(ss.MultiHopOnly, other.MultiHopOnly)

	ss.OpenVPN.mergeWith(other.OpenVPN)
	ss.Wireguard.mergeWith(other.Wireguard)
}

func (ss *ServerSelection) overrideWith(other ServerSelection) {
	ss.VPN = gosettings.OverrideWithString(ss.VPN, other.VPN)
	ss.TargetIP = gosettings.OverrideWithValidator(ss.TargetIP, other.TargetIP)
	ss.Countries = gosettings.OverrideWithSlice(ss.Countries, other.Countries)
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
	ss.OpenVPN.overrideWith(other.OpenVPN)
	ss.Wireguard.overrideWith(other.Wireguard)
}

func (ss *ServerSelection) setDefaults(vpnProvider string) {
	ss.VPN = gosettings.DefaultString(ss.VPN, vpn.OpenVPN)
	ss.TargetIP = gosettings.DefaultValidator(ss.TargetIP, netip.IPv4Unspecified())
	ss.OwnedOnly = gosettings.DefaultPointer(ss.OwnedOnly, false)
	ss.FreeOnly = gosettings.DefaultPointer(ss.FreeOnly, false)
	ss.PremiumOnly = gosettings.DefaultPointer(ss.PremiumOnly, false)
	ss.StreamOnly = gosettings.DefaultPointer(ss.StreamOnly, false)
	ss.MultiHopOnly = gosettings.DefaultPointer(ss.MultiHopOnly, false)
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
