package settings

import (
	"fmt"
	"net"

	"github.com/qdm12/gluetun/internal/configuration/settings/helpers"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
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
	// OwnedOnly is true if only VPN provider owned servers
	// should be filtered. This is used with Mullvad.
	OwnedOnly *bool
	// FreeOnly is true if only free VPN servers
	// should be filtered. This is used with ProtonVPN.
	FreeOnly *bool
	// FreeOnly is true if only free VPN servers
	// should be filtered. This is used with ProtonVPN.
	StreamOnly *bool
	// MultiHopOnly is true if only multihop VPN servers
	// should be filtered. This is used with Surfshark.
	MultiHopOnly *bool

	// OpenVPN contains settings to select OpenVPN servers
	// and the final connection.
	OpenVPN OpenVPNSelection
	// Wireguard contains settings to select Wireguard servers
	// and the final connection.
	Wireguard WireguardSelection
}

func (ss *ServerSelection) validate(vpnServiceProvider string,
	allServers models.AllServers) (err error) {
	switch ss.VPN {
	case constants.OpenVPN, constants.Wireguard:
	default:
		return fmt.Errorf("%w: %s", ErrVPNTypeNotValid, ss.VPN)
	}

	var countryChoices, regionChoices, cityChoices,
		ispChoices, nameChoices, hostnameChoices []string
	switch vpnServiceProvider {
	case constants.Cyberghost:
		servers := allServers.GetCyberghost()
		countryChoices = constants.CyberghostCountryChoices(servers)
		hostnameChoices = constants.CyberghostHostnameChoices(servers)
	case constants.Expressvpn:
		servers := allServers.GetExpressvpn()
		countryChoices = constants.ExpressvpnCountriesChoices(servers)
		cityChoices = constants.ExpressvpnCityChoices(servers)
		hostnameChoices = constants.ExpressvpnHostnameChoices(servers)
	case constants.Fastestvpn:
		servers := allServers.GetFastestvpn()
		countryChoices = constants.FastestvpnCountriesChoices(servers)
		hostnameChoices = constants.FastestvpnHostnameChoices(servers)
	case constants.HideMyAss:
		servers := allServers.GetHideMyAss()
		countryChoices = constants.HideMyAssCountryChoices(servers)
		regionChoices = constants.HideMyAssRegionChoices(servers)
		cityChoices = constants.HideMyAssCityChoices(servers)
		hostnameChoices = constants.HideMyAssHostnameChoices(servers)
	case constants.Ipvanish:
		servers := allServers.GetIpvanish()
		countryChoices = constants.IpvanishCountryChoices(servers)
		cityChoices = constants.IpvanishCityChoices(servers)
		hostnameChoices = constants.IpvanishHostnameChoices(servers)
	case constants.Ivpn:
		servers := allServers.GetIvpn()
		countryChoices = constants.IvpnCountryChoices(servers)
		cityChoices = constants.IvpnCityChoices(servers)
		ispChoices = constants.IvpnISPChoices(servers)
		hostnameChoices = constants.IvpnHostnameChoices(servers)
	case constants.Mullvad:
		servers := allServers.GetMullvad()
		countryChoices = constants.MullvadCountryChoices(servers)
		cityChoices = constants.MullvadCityChoices(servers)
		ispChoices = constants.MullvadISPChoices(servers)
		hostnameChoices = constants.MullvadHostnameChoices(servers)
	case constants.Nordvpn:
		servers := allServers.GetNordvpn()
		regionChoices = constants.NordvpnRegionChoices(servers)
		hostnameChoices = constants.NordvpnHostnameChoices(servers)
	case constants.Perfectprivacy:
		servers := allServers.GetPerfectprivacy()
		cityChoices = constants.PerfectprivacyCityChoices(servers)
	case constants.Privado:
		servers := allServers.GetPrivado()
		countryChoices = constants.PrivadoCountryChoices(servers)
		regionChoices = constants.PrivadoRegionChoices(servers)
		cityChoices = constants.PrivadoCityChoices(servers)
		hostnameChoices = constants.PrivadoHostnameChoices(servers)
	case constants.PrivateInternetAccess:
		servers := allServers.GetPia()
		regionChoices = constants.PIAGeoChoices(servers)
		hostnameChoices = constants.PIAHostnameChoices(servers)
		nameChoices = constants.PIANameChoices(servers)
	case constants.Privatevpn:
		servers := allServers.GetPrivatevpn()
		countryChoices = constants.PrivatevpnCountryChoices(servers)
		cityChoices = constants.PrivatevpnCityChoices(servers)
		hostnameChoices = constants.PrivatevpnHostnameChoices(servers)
	case constants.Protonvpn:
		servers := allServers.GetProtonvpn()
		countryChoices = constants.ProtonvpnCountryChoices(servers)
		regionChoices = constants.ProtonvpnRegionChoices(servers)
		cityChoices = constants.ProtonvpnCityChoices(servers)
		nameChoices = constants.ProtonvpnNameChoices(servers)
		hostnameChoices = constants.ProtonvpnHostnameChoices(servers)
	case constants.Purevpn:
		servers := allServers.GetPurevpn()
		countryChoices = constants.PurevpnCountryChoices(servers)
		regionChoices = constants.PurevpnRegionChoices(servers)
		cityChoices = constants.PurevpnCityChoices(servers)
		hostnameChoices = constants.PurevpnHostnameChoices(servers)
	case constants.Surfshark:
		servers := allServers.GetSurfshark()
		countryChoices = constants.SurfsharkCountryChoices(servers)
		cityChoices = constants.SurfsharkCityChoices(servers)
		hostnameChoices = constants.SurfsharkHostnameChoices(servers)
		regionChoices = constants.SurfsharkRegionChoices(servers)
		// TODO v4 remove
		regionChoices = append(regionChoices, constants.SurfsharkRetroLocChoices(servers)...)
		if err := helpers.AreAllOneOf(ss.Regions, regionChoices); err != nil {
			return fmt.Errorf("%w: %s", ErrRegionNotValid, err)
		}
		// Retro compatibility
		// TODO remove in v4
		*ss = surfsharkRetroRegion(*ss)
	case constants.Torguard:
		servers := allServers.GetTorguard()
		countryChoices = constants.TorguardCountryChoices(servers)
		cityChoices = constants.TorguardCityChoices(servers)
		hostnameChoices = constants.TorguardHostnameChoices(servers)
	case constants.VPNUnlimited:
		servers := allServers.GetVPNUnlimited()
		countryChoices = constants.VPNUnlimitedCountryChoices(servers)
		cityChoices = constants.VPNUnlimitedCityChoices(servers)
		hostnameChoices = constants.VPNUnlimitedHostnameChoices(servers)
	case constants.Vyprvpn:
		servers := allServers.GetVyprvpn()
		regionChoices = constants.VyprvpnRegionChoices(servers)
	case constants.Wevpn:
		servers := allServers.GetWevpn()
		cityChoices = constants.WevpnCityChoices(servers)
		hostnameChoices = constants.WevpnHostnameChoices(servers)
	case constants.Windscribe:
		servers := allServers.GetWindscribe()
		regionChoices = constants.WindscribeRegionChoices(servers)
		cityChoices = constants.WindscribeCityChoices(servers)
		hostnameChoices = constants.WindscribeHostnameChoices(servers)
	default:
		return fmt.Errorf("%w: %s", ErrVPNProviderNameNotValid, ss.VPN)
	}

	err = validateServerFilters(*ss, countryChoices, regionChoices, cityChoices,
		ispChoices, nameChoices, hostnameChoices)
	if err != nil {
		return err // already wrapped error
	}

	err = ss.OpenVPN.validate(vpnServiceProvider)
	if err != nil {
		return fmt.Errorf("OpenVPN server selection settings validation failed: %w", err)
	}

	err = ss.Wireguard.validate(vpnServiceProvider)
	if err != nil {
		return fmt.Errorf("Wireguard server selection settings validation failed: %w", err)
	}

	return nil
}

// validateServerFilters validates filters against the choices given as arguments.
// Set an argument to nil to pass the check for a particular filter.
func validateServerFilters(settings ServerSelection,
	countryChoices, regionChoices, cityChoices, ispChoices,
	nameChoices, hostnameChoices []string) (err error) {
	if countryChoices != nil {
		if err := helpers.AreAllOneOf(settings.Countries, countryChoices); err != nil {
			return fmt.Errorf("%w: %s", ErrCountryNotValid, err)
		}
	}

	if regionChoices != nil {
		if err := helpers.AreAllOneOf(settings.Regions, regionChoices); err != nil {
			return fmt.Errorf("%w: %s", ErrRegionNotValid, err)
		}
	}

	if cityChoices != nil {
		if err := helpers.AreAllOneOf(settings.Cities, cityChoices); err != nil {
			return fmt.Errorf("%w: %s", ErrCityNotValid, err)
		}
	}

	if ispChoices != nil {
		if err := helpers.AreAllOneOf(settings.ISPs, ispChoices); err != nil {
			return fmt.Errorf("%w: %s", ErrISPNotValid, err)
		}
	}

	if hostnameChoices != nil {
		if err := helpers.AreAllOneOf(settings.Hostnames, hostnameChoices); err != nil {
			return fmt.Errorf("%w: %s", ErrHostnameNotValid, err)
		}
	}

	if nameChoices != nil {
		if err := helpers.AreAllOneOf(settings.Names, nameChoices); err != nil {
			return fmt.Errorf("%w: %s", ErrNameNotValid, err)
		}
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
	ss.Names = helpers.MergeStringSlices(ss.Hostnames, other.Names)
	ss.Numbers = helpers.MergeUint16Slices(ss.Numbers, other.Numbers)
	ss.OwnedOnly = helpers.MergeWithBool(ss.OwnedOnly, other.OwnedOnly)
	ss.FreeOnly = helpers.MergeWithBool(ss.FreeOnly, other.FreeOnly)
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
	ss.Names = helpers.OverrideWithStringSlice(ss.Hostnames, other.Names)
	ss.Numbers = helpers.OverrideWithUint16Slice(ss.Numbers, other.Numbers)
	ss.OwnedOnly = helpers.OverrideWithBool(ss.OwnedOnly, other.OwnedOnly)
	ss.FreeOnly = helpers.OverrideWithBool(ss.FreeOnly, other.FreeOnly)
	ss.StreamOnly = helpers.OverrideWithBool(ss.StreamOnly, other.StreamOnly)
	ss.MultiHopOnly = helpers.OverrideWithBool(ss.MultiHopOnly, other.MultiHopOnly)
	ss.OpenVPN.overrideWith(other.OpenVPN)
	ss.Wireguard.overrideWith(other.Wireguard)
}

func (ss *ServerSelection) setDefaults(vpnProvider string) {
	ss.VPN = helpers.DefaultString(ss.VPN, constants.OpenVPN)
	ss.TargetIP = helpers.DefaultIP(ss.TargetIP, net.IP{})
	ss.OwnedOnly = helpers.DefaultBool(ss.OwnedOnly, false)
	ss.FreeOnly = helpers.DefaultBool(ss.FreeOnly, false)
	ss.StreamOnly = helpers.DefaultBool(ss.StreamOnly, false)
	ss.MultiHopOnly = helpers.DefaultBool(ss.MultiHopOnly, false)
	ss.OpenVPN.setDefaults(vpnProvider)
	ss.Wireguard.setDefaults()
}
