package settings

import (
	"fmt"

	"github.com/qdm12/gluetun/internal/configuration/settings/helpers"
	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gluetun/internal/provider/privateinternetaccess/presets"
	"github.com/qdm12/gotree"
)

type OpenVPNSelection struct {
	// ConfFile is the custom configuration file path.
	// It can be set to an empty string to indicate to
	// NOT use a custom configuration file.
	// It cannot be nil in the internal state.
	ConfFile *string
	// TCP is true if the OpenVPN protocol is TCP,
	// and false for UDP.
	// It cannot be nil in the internal state.
	TCP *bool
	// CustomPort is the OpenVPN server endpoint port.
	// It can be set to 0 to indicate no custom port should
	// be used. It cannot be nil in the internal state.
	CustomPort *uint16 // HideMyAss, Mullvad, PIA, ProtonVPN, WeVPN, Windscribe
	// PIAEncPreset is the encryption preset for
	// Private Internet Access. It can be set to an
	// empty string for other providers.
	PIAEncPreset *string
}

func (o OpenVPNSelection) validate(vpnProvider string) (err error) {
	// Validate ConfFile
	if confFile := *o.ConfFile; confFile != "" {
		err := helpers.FileExists(confFile)
		if err != nil {
			return fmt.Errorf("configuration file: %w", err)
		}
	}

	// Validate TCP
	if *o.TCP && helpers.IsOneOf(vpnProvider,
		providers.Ipvanish,
		providers.Perfectprivacy,
		providers.Privado,
		providers.VPNUnlimited,
		providers.Vyprvpn,
	) {
		return fmt.Errorf("%w: for VPN service provider %s",
			ErrOpenVPNTCPNotSupported, vpnProvider)
	}

	// Validate CustomPort
	if *o.CustomPort != 0 {
		switch vpnProvider {
		// no restriction on port
		case providers.Cyberghost, providers.HideMyAss,
			providers.Privatevpn, providers.Torguard:
		// no custom port allowed
		case providers.Expressvpn, providers.Fastestvpn,
			providers.Ipvanish, providers.Nordvpn,
			providers.Privado, providers.Purevpn,
			providers.Surfshark, providers.VPNSecure,
			providers.VPNUnlimited, providers.Vyprvpn:
			return fmt.Errorf("%w: for VPN service provider %s",
				ErrOpenVPNCustomPortNotAllowed, vpnProvider)
		default:
			var allowedTCP, allowedUDP []uint16
			switch vpnProvider {
			case providers.Airvpn:
				allowedTCP = []uint16{
					53, 80, 443, // IP in 1, 3
					1194, 2018, 41185, // IP in 1, 2, 3, 4
				}
				allowedUDP = []uint16{53, 80, 443, 1194, 2018, 41185}
			case providers.Ivpn:
				allowedTCP = []uint16{80, 443, 1143}
				allowedUDP = []uint16{53, 1194, 2049, 2050}
			case providers.Mullvad:
				allowedTCP = []uint16{80, 443, 1401}
				allowedUDP = []uint16{53, 1194, 1195, 1196, 1197, 1300, 1301, 1302, 1303, 1400}
			case providers.Perfectprivacy:
				allowedTCP = []uint16{44, 443, 4433}
				allowedUDP = []uint16{44, 443, 4433}
			case providers.PrivateInternetAccess:
				allowedTCP = []uint16{80, 110, 443}
				allowedUDP = []uint16{53, 1194, 1197, 1198, 8080, 9201}
			case providers.Protonvpn:
				allowedTCP = []uint16{443, 5995, 8443}
				allowedUDP = []uint16{80, 443, 1194, 4569, 5060}
			case providers.SlickVPN:
				allowedTCP = []uint16{443, 8080, 8888}
				allowedUDP = []uint16{443, 8080, 8888}
			case providers.Wevpn:
				allowedTCP = []uint16{53, 1195, 1199, 2018}
				allowedUDP = []uint16{80, 1194, 1198}
			case providers.Windscribe:
				allowedTCP = []uint16{21, 22, 80, 123, 143, 443, 587, 1194, 3306, 8080, 54783}
				allowedUDP = []uint16{53, 80, 123, 443, 1194, 54783}
			}

			if *o.TCP && !helpers.Uint16IsOneOf(*o.CustomPort, allowedTCP) {
				return fmt.Errorf("%w: %d for VPN service provider %s; %s",
					ErrOpenVPNCustomPortNotAllowed, o.CustomPort, vpnProvider,
					helpers.PortChoicesOrString(allowedTCP))
			} else if !*o.TCP && !helpers.Uint16IsOneOf(*o.CustomPort, allowedUDP) {
				return fmt.Errorf("%w: %d for VPN service provider %s; %s",
					ErrOpenVPNCustomPortNotAllowed, o.CustomPort, vpnProvider,
					helpers.PortChoicesOrString(allowedUDP))
			}
		}
	}

	// Validate EncPreset
	if vpnProvider == providers.PrivateInternetAccess {
		validEncryptionPresets := []string{
			presets.None,
			presets.Normal,
			presets.Strong,
		}
		if !helpers.IsOneOf(*o.PIAEncPreset, validEncryptionPresets...) {
			return fmt.Errorf("%w: %s; valid presets are %s",
				ErrOpenVPNEncryptionPresetNotValid, *o.PIAEncPreset,
				helpers.ChoicesOrString(validEncryptionPresets))
		}
	}

	return nil
}

func (o *OpenVPNSelection) copy() (copied OpenVPNSelection) {
	return OpenVPNSelection{
		ConfFile:     helpers.CopyStringPtr(o.ConfFile),
		TCP:          helpers.CopyBoolPtr(o.TCP),
		CustomPort:   helpers.CopyUint16Ptr(o.CustomPort),
		PIAEncPreset: helpers.CopyStringPtr(o.PIAEncPreset),
	}
}

func (o *OpenVPNSelection) mergeWith(other OpenVPNSelection) {
	o.ConfFile = helpers.MergeWithStringPtr(o.ConfFile, other.ConfFile)
	o.TCP = helpers.MergeWithBool(o.TCP, other.TCP)
	o.CustomPort = helpers.MergeWithUint16(o.CustomPort, other.CustomPort)
	o.PIAEncPreset = helpers.MergeWithStringPtr(o.PIAEncPreset, other.PIAEncPreset)
}

func (o *OpenVPNSelection) overrideWith(other OpenVPNSelection) {
	o.ConfFile = helpers.OverrideWithStringPtr(o.ConfFile, other.ConfFile)
	o.TCP = helpers.OverrideWithBool(o.TCP, other.TCP)
	o.CustomPort = helpers.OverrideWithUint16(o.CustomPort, other.CustomPort)
	o.PIAEncPreset = helpers.OverrideWithStringPtr(o.PIAEncPreset, other.PIAEncPreset)
}

func (o *OpenVPNSelection) setDefaults(vpnProvider string) {
	o.ConfFile = helpers.DefaultStringPtr(o.ConfFile, "")
	o.TCP = helpers.DefaultBool(o.TCP, false)
	o.CustomPort = helpers.DefaultUint16(o.CustomPort, 0)

	var defaultEncPreset string
	if vpnProvider == providers.PrivateInternetAccess {
		defaultEncPreset = presets.Strong
	}
	o.PIAEncPreset = helpers.DefaultStringPtr(o.PIAEncPreset, defaultEncPreset)
}

func (o OpenVPNSelection) String() string {
	return o.toLinesNode().String()
}

func (o OpenVPNSelection) toLinesNode() (node *gotree.Node) {
	node = gotree.New("OpenVPN server selection settings:")
	node.Appendf("Protocol: %s", helpers.TCPPtrToString(o.TCP))

	if *o.CustomPort != 0 {
		node.Appendf("Custom port: %d", *o.CustomPort)
	}

	if *o.PIAEncPreset != "" {
		node.Appendf("Private Internet Access encryption preset: %s", *o.PIAEncPreset)
	}

	if *o.ConfFile != "" {
		node.Appendf("Custom configuration file: %s", *o.ConfFile)
	}

	return node
}
