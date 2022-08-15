package settings

import (
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/qdm12/gluetun/internal/configuration/settings/helpers"
	"github.com/qdm12/gluetun/internal/configuration/settings/validation"
	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gluetun/internal/constants/vpn"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gotree"
)

type ServerSelection struct { //nolint:maligned
	// VPN is the VPN type which can be 'openvpn'
	// or 'wireguard'. It cannot be the empty string
	// in the internal state.
	VPN string
	// TargetIP is the server endpoint IP address to use.
	// It will override any IP address from the picked
	// built-in server. It cannot be nil in the internal
	// state, and can be set to an empty net.IP{} to indicate
	// there is not target IP address to use.
	TargetIP net.IP
	// Counties is the list of countries to filter VPN servers with.
	Countries []string
	// Regions is the list of regions to filter VPN servers with.
	Regions []string
	// Cities is the list of cities to filter VPN servers with.
	Cities []string
	// ISPs is the list of ISP names to filter VPN servers with.
	ISPs []string
	// Names is the list of server names to filter VPN servers with.
	Names []string
	// Numbers is the list of server numbers to filter VPN servers with.
	Numbers []uint16
	// Hostnames is the list of hostnames to filter VPN servers with.
	Hostnames []string
	// OwnedOnly is true if VPN provider servers that are not owned
	// should be filtered. This is used with Mullvad.
	OwnedOnly *bool
	// FreeOnly is true if VPN servers that are not free should
	// be filtered. This is used with ProtonVPN and VPN Unlimited.
	FreeOnly *bool
	// PremiumOnly is true if VPN servers that are not premium should
	// be filtered. This is used with VPN Secure.
	// TODO extend to providers using FreeOnly.
	PremiumOnly *bool
	// StreamOnly is true if VPN servers not for streaming should
	// be filtered. This is used with VPNUnlimited.
	StreamOnly *bool
	// MultiHopOnly is true if VPN servers that are not multihop
	// should be filtered. This is used with Surfshark.
	MultiHopOnly *bool

	// OpenVPN contains settings to select OpenVPN servers
	// and the final connection.
	OpenVPN OpenVPNSelection
	// Wireguard contains settings to select Wireguard servers
	// and the final connection.
	Wireguard WireguardSelection
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

	err = validateServerFilters(*ss, filterChoices)
	if err != nil {
		if errors.Is(err, helpers.ErrNoChoice) {
			return fmt.Errorf("for VPN service provider %s: %w", vpnServiceProvider, err)
		}
		return err // already wrapped error
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
		return ErrFreePremiumBothSet
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
		if err := helpers.AreAllOneOf(ss.Regions, filterChoices.Regions); err != nil {
			return models.FilterChoices{}, fmt.Errorf("%w: %s", ErrRegionNotValid, err)
		}
		*ss = surfsharkRetroRegion(*ss)
	}

	return filterChoices, nil
}

// validateServerFilters validates filters against the choices given as arguments.
// Set an argument to nil to pass the check for a particular filter.
func validateServerFilters(settings ServerSelection, filterChoices models.FilterChoices) (err error) {
	if err := helpers.AreAllOneOf(settings.Countries, filterChoices.Countries); err != nil {
		return fmt.Errorf("%w: %s", ErrCountryNotValid, err)
	}

	if err := helpers.AreAllOneOf(settings.Regions, filterChoices.Regions); err != nil {
		return fmt.Errorf("%w: %s", ErrRegionNotValid, err)
	}

	if err := helpers.AreAllOneOf(settings.Cities, filterChoices.Cities); err != nil {
		return fmt.Errorf("%w: %s", ErrCityNotValid, err)
	}

	if err := helpers.AreAllOneOf(settings.ISPs, filterChoices.ISPs); err != nil {
		return fmt.Errorf("%w: %s", ErrISPNotValid, err)
	}

	if err := helpers.AreAllOneOf(settings.Hostnames, filterChoices.Hostnames); err != nil {
		return fmt.Errorf("%w: %s", ErrHostnameNotValid, err)
	}

	if err := helpers.AreAllOneOf(settings.Names, filterChoices.Names); err != nil {
		return fmt.Errorf("%w: %s", ErrNameNotValid, err)
	}

	return nil
}

func (ss *ServerSelection) copy() (copied ServerSelection) {
	return ServerSelection{
		VPN:          ss.VPN,
		TargetIP:     helpers.CopyIP(ss.TargetIP),
		Countries:    helpers.CopyStringSlice(ss.Countries),
		Regions:      helpers.CopyStringSlice(ss.Regions),
		Cities:       helpers.CopyStringSlice(ss.Cities),
		ISPs:         helpers.CopyStringSlice(ss.ISPs),
		Hostnames:    helpers.CopyStringSlice(ss.Hostnames),
		Names:        helpers.CopyStringSlice(ss.Names),
		Numbers:      helpers.CopyUint16Slice(ss.Numbers),
		OwnedOnly:    helpers.CopyBoolPtr(ss.OwnedOnly),
		FreeOnly:     helpers.CopyBoolPtr(ss.FreeOnly),
		PremiumOnly:  helpers.CopyBoolPtr(ss.PremiumOnly),
		StreamOnly:   helpers.CopyBoolPtr(ss.StreamOnly),
		MultiHopOnly: helpers.CopyBoolPtr(ss.MultiHopOnly),
		OpenVPN:      ss.OpenVPN.copy(),
		Wireguard:    ss.Wireguard.copy(),
	}
}

func (ss *ServerSelection) mergeWith(other ServerSelection) {
	ss.VPN = helpers.MergeWithString(ss.VPN, other.VPN)
	ss.TargetIP = helpers.MergeWithIP(ss.TargetIP, other.TargetIP)
	ss.Countries = helpers.MergeStringSlices(ss.Countries, other.Countries)
	ss.Regions = helpers.MergeStringSlices(ss.Regions, other.Regions)
	ss.Cities = helpers.MergeStringSlices(ss.Cities, other.Cities)
	ss.ISPs = helpers.MergeStringSlices(ss.ISPs, other.ISPs)
	ss.Hostnames = helpers.MergeStringSlices(ss.Hostnames, other.Hostnames)
	ss.Names = helpers.MergeStringSlices(ss.Names, other.Names)
	ss.Numbers = helpers.MergeUint16Slices(ss.Numbers, other.Numbers)
	ss.OwnedOnly = helpers.MergeWithBool(ss.OwnedOnly, other.OwnedOnly)
	ss.FreeOnly = helpers.MergeWithBool(ss.FreeOnly, other.FreeOnly)
	ss.PremiumOnly = helpers.MergeWithBool(ss.PremiumOnly, other.PremiumOnly)
	ss.StreamOnly = helpers.MergeWithBool(ss.StreamOnly, other.StreamOnly)
	ss.MultiHopOnly = helpers.MergeWithBool(ss.MultiHopOnly, other.MultiHopOnly)

	ss.OpenVPN.mergeWith(other.OpenVPN)
	ss.Wireguard.mergeWith(other.Wireguard)
}

func (ss *ServerSelection) overrideWith(other ServerSelection) {
	ss.VPN = helpers.OverrideWithString(ss.VPN, other.VPN)
	ss.TargetIP = helpers.OverrideWithIP(ss.TargetIP, other.TargetIP)
	ss.Countries = helpers.OverrideWithStringSlice(ss.Countries, other.Countries)
	ss.Regions = helpers.OverrideWithStringSlice(ss.Regions, other.Regions)
	ss.Cities = helpers.OverrideWithStringSlice(ss.Cities, other.Cities)
	ss.ISPs = helpers.OverrideWithStringSlice(ss.ISPs, other.ISPs)
	ss.Hostnames = helpers.OverrideWithStringSlice(ss.Hostnames, other.Hostnames)
	ss.Names = helpers.OverrideWithStringSlice(ss.Names, other.Names)
	ss.Numbers = helpers.OverrideWithUint16Slice(ss.Numbers, other.Numbers)
	ss.OwnedOnly = helpers.OverrideWithBool(ss.OwnedOnly, other.OwnedOnly)
	ss.FreeOnly = helpers.OverrideWithBool(ss.FreeOnly, other.FreeOnly)
	ss.PremiumOnly = helpers.OverrideWithBool(ss.PremiumOnly, other.PremiumOnly)
	ss.StreamOnly = helpers.OverrideWithBool(ss.StreamOnly, other.StreamOnly)
	ss.MultiHopOnly = helpers.OverrideWithBool(ss.MultiHopOnly, other.MultiHopOnly)
	ss.OpenVPN.overrideWith(other.OpenVPN)
	ss.Wireguard.overrideWith(other.Wireguard)
}

func (ss *ServerSelection) setDefaults(vpnProvider string) {
	ss.VPN = helpers.DefaultString(ss.VPN, vpn.OpenVPN)
	ss.TargetIP = helpers.DefaultIP(ss.TargetIP, net.IP{})
	ss.OwnedOnly = helpers.DefaultBool(ss.OwnedOnly, false)
	ss.FreeOnly = helpers.DefaultBool(ss.FreeOnly, false)
	ss.PremiumOnly = helpers.DefaultBool(ss.PremiumOnly, false)
	ss.StreamOnly = helpers.DefaultBool(ss.StreamOnly, false)
	ss.MultiHopOnly = helpers.DefaultBool(ss.MultiHopOnly, false)
	ss.OpenVPN.setDefaults(vpnProvider)
	ss.Wireguard.setDefaults()
}

func (ss ServerSelection) String() string {
	return ss.toLinesNode().String()
}

func (ss ServerSelection) toLinesNode() (node *gotree.Node) {
	node = gotree.New("Server selection settings:")
	node.Appendf("VPN type: %s", ss.VPN)
	if len(ss.TargetIP) > 0 {
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
