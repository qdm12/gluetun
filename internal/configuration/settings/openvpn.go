package settings

import (
	"fmt"
	"strings"

	"github.com/qdm12/gluetun/internal/configuration/settings/helpers"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/openvpn/parse"
	"github.com/qdm12/gotree"
)

// OpenVPN contains settings to configure the OpenVPN client.
type OpenVPN struct {
	// Version is the OpenVPN version to run.
	// It can only be "2.4" or "2.5".
	Version string
	// User is the OpenVPN authentication username.
	// It cannot be an empty string in the internal state
	// if OpenVPN is used.
	User string
	// Password is the OpenVPN authentication password.
	// It cannot be an empty string in the internal state
	// if OpenVPN is used.
	Password string
	// ConfFile is a custom OpenVPN configuration file path.
	// It can be set to the empty string for it to be ignored.
	// It cannot be nil in the internal state.
	ConfFile *string
	// Ciphers is a list of ciphers to use for OpenVPN,
	// different from the ones specified by the VPN
	// service provider configuration files.
	Ciphers []string
	// Auth is an auth algorithm to use in OpenVPN instead
	// of the one specified by the VPN service provider.
	// It cannot be nil in the internal state.
	// It is ignored if it is set to the empty string.
	Auth *string
	// ClientCrt is the OpenVPN client certificate.
	// This is notably used by Cyberghost.
	// It can be set to the empty string to be ignored.
	// It cannot be nil in the internal state.
	ClientCrt *string
	// ClientKey is the OpenVPN client key.
	// This is used by Cyberghost and VPN Unlimited.
	// It can be set to the empty string to be ignored.
	// It cannot be nil in the internal state.
	ClientKey *string
	// PIAEncPreset is the encryption preset for
	// Private Internet Access. It can be set to an
	// empty string for other providers.
	PIAEncPreset *string
	// IPv6 is set to true if IPv6 routing should be
	// set to be tunnel in OpenVPN, and false otherwise.
	// It cannot be nil in the internal state.
	IPv6 *bool // TODO automate like with Wireguard
	// MSSFix is the value (1 to 10000) to set for the
	// mssfix option for OpenVPN. It is ignored if set to 0.
	// It cannot be nil in the internal state.
	MSSFix *uint16
	// Interface is the OpenVPN device interface name.
	// It cannot be an empty string in the internal state.
	Interface string
	// Root is true if OpenVPN is to be run as root,
	// and false otherwise. It cannot be nil in the
	// internal state.
	Root *bool
	// ProcUser is the OpenVPN process OS username
	// to use. It cannot be nil in the internal state.
	// This is set and injected at runtime.
	// TODO only use ProcUser and not Root field.
	ProcUser string
	// Verbosity is the OpenVPN verbosity level from 0 to 6.
	// It cannot be nil in the internal state.
	Verbosity *int
	// Flags is a slice of additional flags to be passed
	// to the OpenVPN program.
	Flags []string
}

func (o OpenVPN) validate(vpnProvider string) (err error) {
	// Validate version
	validVersions := []string{constants.Openvpn24, constants.Openvpn25}
	if !helpers.IsOneOf(o.Version, validVersions...) {
		return fmt.Errorf("%w: %q can only be one of %s",
			ErrOpenVPNVersionIsNotValid, o.Version, strings.Join(validVersions, ", "))
	}

	if o.User == "" {
		return ErrOpenVPNUserIsEmpty
	}

	if o.Password == "" {
		return ErrOpenVPNPasswordIsEmpty
	}

	// Validate ConfFile
	if vpnProvider == constants.Custom {
		if *o.ConfFile == "" {
			return fmt.Errorf("%w: no file path specified", ErrOpenVPNConfigFile)
		}
		err := helpers.FileExists(*o.ConfFile)
		if err != nil {
			return fmt.Errorf("%w: %s", ErrOpenVPNConfigFile, err)
		}
	}

	// Check client certificate
	switch vpnProvider {
	case
		constants.Cyberghost,
		constants.VPNUnlimited:
		if *o.ClientCrt == "" {
			return ErrOpenVPNClientCertMissing
		}
	}
	if *o.ClientCrt != "" {
		_, err = parse.ExtractCert([]byte(*o.ClientCrt))
		if err != nil {
			return fmt.Errorf("%w: %s", ErrOpenVPNClientCertNotValid, err)
		}
	}

	// Check client key
	switch vpnProvider {
	case
		constants.Cyberghost,
		constants.VPNUnlimited,
		constants.Wevpn:
		if *o.ClientKey == "" {
			return ErrOpenVPNClientKeyMissing
		}
	}
	if *o.ClientKey != "" {
		_, err = parse.ExtractPrivateKey([]byte(*o.ClientKey))
		if err != nil {
			return fmt.Errorf("%w: %s", ErrOpenVPNClientKeyNotValid, err)
		}
	}

	// Validate MSSFix
	const maxMSSFix = 10000
	if *o.MSSFix > maxMSSFix {
		return fmt.Errorf("%w: %d is over the maximum value of %d",
			ErrOpenVPNMSSFixIsTooHigh, *o.MSSFix, maxMSSFix)
	}

	if !regexpInterfaceName.MatchString(o.Interface) {
		return fmt.Errorf("%w: '%s' does not match regex '%s'",
			ErrOpenVPNInterfaceNotValid, o.Interface, regexpInterfaceName)
	}

	// Validate Verbosity
	if *o.Verbosity < 0 || *o.Verbosity > 6 {
		return fmt.Errorf("%w: %d can only be between 0 and 5",
			ErrOpenVPNVerbosityIsOutOfBounds, o.Verbosity)
	}

	return nil
}

func (o *OpenVPN) copy() (copied OpenVPN) {
	return OpenVPN{
		Version:      o.Version,
		User:         o.User,
		Password:     o.Password,
		ConfFile:     helpers.CopyStringPtr(o.ConfFile),
		Ciphers:      helpers.CopyStringSlice(o.Ciphers),
		Auth:         helpers.CopyStringPtr(o.Auth),
		ClientCrt:    helpers.CopyStringPtr(o.ClientCrt),
		ClientKey:    helpers.CopyStringPtr(o.ClientKey),
		PIAEncPreset: helpers.CopyStringPtr(o.PIAEncPreset),
		IPv6:         helpers.CopyBoolPtr(o.IPv6),
		MSSFix:       helpers.CopyUint16Ptr(o.MSSFix),
		Interface:    o.Interface,
		Root:         helpers.CopyBoolPtr(o.Root),
		ProcUser:     o.ProcUser,
		Verbosity:    helpers.CopyIntPtr(o.Verbosity),
		Flags:        helpers.CopyStringSlice(o.Flags),
	}
}

// mergeWith merges the other settings into any
// unset field of the receiver settings object.
func (o *OpenVPN) mergeWith(other OpenVPN) {
	o.Version = helpers.MergeWithString(o.Version, other.Version)
	o.User = helpers.MergeWithString(o.User, other.User)
	o.Password = helpers.MergeWithString(o.Password, other.Password)
	o.ConfFile = helpers.MergeWithStringPtr(o.ConfFile, other.ConfFile)
	o.Ciphers = helpers.MergeStringSlices(o.Ciphers, other.Ciphers)
	o.Auth = helpers.MergeWithStringPtr(o.Auth, other.Auth)
	o.ClientCrt = helpers.MergeWithStringPtr(o.ClientCrt, other.ClientCrt)
	o.ClientKey = helpers.MergeWithStringPtr(o.ClientKey, other.ClientKey)
	o.PIAEncPreset = helpers.MergeWithStringPtr(o.PIAEncPreset, other.PIAEncPreset)
	o.IPv6 = helpers.MergeWithBool(o.IPv6, other.IPv6)
	o.MSSFix = helpers.MergeWithUint16(o.MSSFix, other.MSSFix)
	o.Interface = helpers.MergeWithString(o.Interface, other.Interface)
	o.Root = helpers.MergeWithBool(o.Root, other.Root)
	o.ProcUser = helpers.MergeWithString(o.ProcUser, other.ProcUser)
	o.Verbosity = helpers.MergeWithInt(o.Verbosity, other.Verbosity)
	o.Flags = helpers.MergeStringSlices(o.Flags, other.Flags)
}

// overrideWith overrides fields of the receiver
// settings object with any field set in the other
// settings.
func (o *OpenVPN) overrideWith(other OpenVPN) {
	o.Version = helpers.OverrideWithString(o.Version, other.Version)
	o.User = helpers.OverrideWithString(o.User, other.User)
	o.Password = helpers.OverrideWithString(o.Password, other.Password)
	o.ConfFile = helpers.OverrideWithStringPtr(o.ConfFile, other.ConfFile)
	o.Ciphers = helpers.OverrideWithStringSlice(o.Ciphers, other.Ciphers)
	o.Auth = helpers.OverrideWithStringPtr(o.Auth, other.Auth)
	o.ClientCrt = helpers.OverrideWithStringPtr(o.ClientCrt, other.ClientCrt)
	o.ClientKey = helpers.OverrideWithStringPtr(o.ClientKey, other.ClientKey)
	o.PIAEncPreset = helpers.OverrideWithStringPtr(o.PIAEncPreset, other.PIAEncPreset)
	o.IPv6 = helpers.OverrideWithBool(o.IPv6, other.IPv6)
	o.MSSFix = helpers.OverrideWithUint16(o.MSSFix, other.MSSFix)
	o.Interface = helpers.OverrideWithString(o.Interface, other.Interface)
	o.Root = helpers.OverrideWithBool(o.Root, other.Root)
	o.ProcUser = helpers.OverrideWithString(o.ProcUser, other.ProcUser)
	o.Verbosity = helpers.OverrideWithInt(o.Verbosity, other.Verbosity)
	o.Flags = helpers.OverrideWithStringSlice(o.Flags, other.Flags)
}

func (o *OpenVPN) setDefaults(vpnProvider string) {
	o.Version = helpers.DefaultString(o.Version, constants.Openvpn25)
	if vpnProvider == constants.Mullvad {
		o.Password = "m"
	}

	o.ConfFile = helpers.DefaultStringPtr(o.ConfFile, "")
	o.Auth = helpers.DefaultStringPtr(o.Auth, "")
	o.ClientCrt = helpers.DefaultStringPtr(o.ClientCrt, "")
	o.ClientKey = helpers.DefaultStringPtr(o.ClientKey, "")

	var defaultEncPreset string
	if vpnProvider == constants.PrivateInternetAccess {
		defaultEncPreset = constants.PIAEncryptionPresetStrong
	}
	o.PIAEncPreset = helpers.DefaultStringPtr(o.PIAEncPreset, defaultEncPreset)

	o.IPv6 = helpers.DefaultBool(o.IPv6, false)
	o.MSSFix = helpers.DefaultUint16(o.MSSFix, 0)
	o.Interface = helpers.DefaultString(o.Interface, "tun0")
	o.Root = helpers.DefaultBool(o.Root, true)
	o.ProcUser = helpers.DefaultString(o.ProcUser, "root")
	o.Verbosity = helpers.DefaultInt(o.Verbosity, 1)
}

func (o OpenVPN) String() string {
	return o.toLinesNode().String()
}

func (o OpenVPN) toLinesNode() (node *gotree.Node) {
	node = gotree.New("OpenVPN server selection settings:")
	node.Appendf("OpenVPN version: %s", o.Version)
	node.Appendf("User: %s", helpers.ObfuscatePassword(o.User))
	node.Appendf("Password: %s", helpers.ObfuscatePassword(o.Password))

	if *o.ConfFile != "" {
		node.Appendf("Custom configuration file: %s", *o.ConfFile)
	}

	if len(o.Ciphers) > 0 {
		node.Appendf("Ciphers: %s", o.Ciphers)
	}

	if *o.Auth != "" {
		node.Appendf("Auth: %s", *o.Auth)
	}

	if *o.ClientCrt != "" {
		node.Appendf("Client crt: %s", helpers.ObfuscateData(*o.ClientCrt))
	}

	if *o.ClientKey != "" {
		node.Appendf("Client key: %s", helpers.ObfuscateData(*o.ClientKey))
	}

	if *o.PIAEncPreset != "" {
		node.Appendf("Private Internet Access encryption preset: %s", *o.PIAEncPreset)
	}

	node.Appendf("Tunnel IPv6: %s", helpers.BoolPtrToYesNo(o.IPv6))

	if *o.MSSFix > 0 {
		node.Appendf("MSS Fix: %d", *o.MSSFix)
	}

	if o.Interface != "" {
		node.Appendf("Network interface: %s", o.Interface)
	}

	processUser := "root"
	if !*o.Root {
		processUser = "some non root user" // TODO
		if o.ProcUser != "" {
			processUser = o.ProcUser
		}
	}
	node.Appendf("Run OpenVPN as: %s", processUser)

	node.Appendf("Verbosity level: %d", *o.Verbosity)

	if len(o.Flags) > 0 {
		node.Appendf("Flags: %s", o.Flags)
	}

	return node
}

// WithDefaults is a shorthand using setDefaults.
// It's used in unit tests in other packages.
func (o OpenVPN) WithDefaults(provider string) OpenVPN {
	o.setDefaults(provider)
	return o
}
