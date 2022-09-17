package settings

import (
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"

	"github.com/qdm12/gluetun/internal/configuration/settings/helpers"
	"github.com/qdm12/gluetun/internal/constants/openvpn"
	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gluetun/internal/openvpn/extract"
	"github.com/qdm12/gluetun/internal/provider/privateinternetaccess/presets"
	"github.com/qdm12/gotree"
)

// OpenVPN contains settings to configure the OpenVPN client.
type OpenVPN struct {
	// Version is the OpenVPN version to run.
	// It can only be "2.4" or "2.5".
	Version string
	// User is the OpenVPN authentication username.
	// It cannot be nil in the internal state if OpenVPN is used.
	// It is usually required but in some cases can be the empty string
	// to indicate no user+password authentication is needed.
	User *string
	// Password is the OpenVPN authentication password.
	// It cannot be nil in the internal state if OpenVPN is used.
	// It is usually required but in some cases can be the empty string
	// to indicate no user+password authentication is needed.
	Password *string
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
	// Cert is the OpenVPN certificate for the <cert> block.
	// This is notably used by Cyberghost and VPN secure.
	// It can be set to the empty string to be ignored.
	// It cannot be nil in the internal state.
	Cert *string
	// Key is the OpenVPN key.
	// This is used by Cyberghost and VPN Unlimited.
	// It can be set to the empty string to be ignored.
	// It cannot be nil in the internal state.
	Key *string
	// EncryptedKey is the content of an encrypted
	// key for OpenVPN. It is used by VPN secure.
	// It defaults to the empty string meaning it is not
	// to be used. KeyPassphrase must be set if this one is set.
	EncryptedKey *string
	// KeyPassphrase is the key passphrase to be used by OpenVPN
	// to decrypt the EncryptedPrivateKey. It defaults to the
	// empty string and must be set if EncryptedPrivateKey is set.
	KeyPassphrase *string
	// PIAEncPreset is the encryption preset for
	// Private Internet Access. It can be set to an
	// empty string for other providers.
	PIAEncPreset *string
	// MSSFix is the value (1 to 10000) to set for the
	// mssfix option for OpenVPN. It is ignored if set to 0.
	// It cannot be nil in the internal state.
	MSSFix *uint16
	// Interface is the OpenVPN device interface name.
	// It cannot be an empty string in the internal state.
	Interface string
	// ProcessUser is the OpenVPN process OS username
	// to use. It cannot be empty in the internal state.
	// It defaults to 'root'.
	ProcessUser string
	// Verbosity is the OpenVPN verbosity level from 0 to 6.
	// It cannot be nil in the internal state.
	Verbosity *int
	// Flags is a slice of additional flags to be passed
	// to the OpenVPN program.
	Flags []string
}

var ivpnAccountID = regexp.MustCompile(`^(i|ivpn)\-[a-zA-Z0-9]{4}\-[a-zA-Z0-9]{4}\-[a-zA-Z0-9]{4}$`)

func (o OpenVPN) validate(vpnProvider string) (err error) {
	// Validate version
	validVersions := []string{openvpn.Openvpn24, openvpn.Openvpn25}
	if !helpers.IsOneOf(o.Version, validVersions...) {
		return fmt.Errorf("%w: %q can only be one of %s",
			ErrOpenVPNVersionIsNotValid, o.Version, strings.Join(validVersions, ", "))
	}

	isCustom := vpnProvider == providers.Custom
	isUserRequired := !isCustom &&
		vpnProvider != providers.Airvpn &&
		vpnProvider != providers.VPNSecure

	if isUserRequired && *o.User == "" {
		return ErrOpenVPNUserIsEmpty
	}

	passwordRequired := isUserRequired &&
		(vpnProvider != providers.Ivpn || !ivpnAccountID.MatchString(*o.User))

	if passwordRequired && *o.Password == "" {
		return ErrOpenVPNPasswordIsEmpty
	}

	err = validateOpenVPNConfigFilepath(isCustom, *o.ConfFile)
	if err != nil {
		return fmt.Errorf("custom configuration file: %w", err)
	}

	err = validateOpenVPNClientCertificate(vpnProvider, *o.Cert)
	if err != nil {
		return fmt.Errorf("client certificate: %w", err)
	}

	err = validateOpenVPNClientKey(vpnProvider, *o.Key)
	if err != nil {
		return fmt.Errorf("client key: %w", err)
	}

	err = validateOpenVPNEncryptedKey(vpnProvider, *o.EncryptedKey)
	if err != nil {
		return fmt.Errorf("encrypted key: %w", err)
	}

	if *o.EncryptedKey != "" && *o.KeyPassphrase == "" {
		return fmt.Errorf("%w", ErrOpenVPNKeyPassphraseIsEmpty)
	}

	const maxMSSFix = 10000
	if *o.MSSFix > maxMSSFix {
		return fmt.Errorf("%w: %d is over the maximum value of %d",
			ErrOpenVPNMSSFixIsTooHigh, *o.MSSFix, maxMSSFix)
	}

	if !regexpInterfaceName.MatchString(o.Interface) {
		return fmt.Errorf("%w: '%s' does not match regex '%s'",
			ErrOpenVPNInterfaceNotValid, o.Interface, regexpInterfaceName)
	}

	if *o.Verbosity < 0 || *o.Verbosity > 6 {
		return fmt.Errorf("%w: %d can only be between 0 and 5",
			ErrOpenVPNVerbosityIsOutOfBounds, o.Verbosity)
	}

	return nil
}

func validateOpenVPNConfigFilepath(isCustom bool,
	confFile string) (err error) {
	if !isCustom {
		return nil
	}

	if confFile == "" {
		return ErrFilepathMissing
	}

	err = helpers.FileExists(confFile)
	if err != nil {
		return err
	}

	extractor := extract.New()
	_, _, err = extractor.Data(confFile)
	if err != nil {
		return fmt.Errorf("failed extracting information from custom configuration file: %w", err)
	}

	return nil
}

func validateOpenVPNClientCertificate(vpnProvider,
	clientCert string) (err error) {
	switch vpnProvider {
	case
		providers.Airvpn,
		providers.Cyberghost,
		providers.VPNSecure,
		providers.VPNUnlimited:
		if clientCert == "" {
			return ErrMissingValue
		}
	}

	if clientCert == "" {
		return nil
	}

	_, err = base64.StdEncoding.DecodeString(clientCert)
	if err != nil {
		return err
	}
	return nil
}

func validateOpenVPNClientKey(vpnProvider, clientKey string) (err error) {
	switch vpnProvider {
	case
		providers.Airvpn,
		providers.Cyberghost,
		providers.VPNUnlimited,
		providers.Wevpn:
		if clientKey == "" {
			return ErrMissingValue
		}
	}

	if clientKey == "" {
		return nil
	}

	_, err = base64.StdEncoding.DecodeString(clientKey)
	if err != nil {
		return err
	}
	return nil
}

func validateOpenVPNEncryptedKey(vpnProvider,
	encryptedPrivateKey string) (err error) {
	if vpnProvider == providers.VPNSecure && encryptedPrivateKey == "" {
		return ErrMissingValue
	}

	if encryptedPrivateKey == "" {
		return nil
	}

	_, err = base64.StdEncoding.DecodeString(encryptedPrivateKey)
	if err != nil {
		return err
	}
	return nil
}

func (o *OpenVPN) copy() (copied OpenVPN) {
	return OpenVPN{
		Version:       o.Version,
		User:          helpers.CopyStringPtr(o.User),
		Password:      helpers.CopyStringPtr(o.Password),
		ConfFile:      helpers.CopyStringPtr(o.ConfFile),
		Ciphers:       helpers.CopyStringSlice(o.Ciphers),
		Auth:          helpers.CopyStringPtr(o.Auth),
		Cert:          helpers.CopyStringPtr(o.Cert),
		Key:           helpers.CopyStringPtr(o.Key),
		EncryptedKey:  helpers.CopyStringPtr(o.EncryptedKey),
		KeyPassphrase: helpers.CopyStringPtr(o.KeyPassphrase),
		PIAEncPreset:  helpers.CopyStringPtr(o.PIAEncPreset),
		MSSFix:        helpers.CopyUint16Ptr(o.MSSFix),
		Interface:     o.Interface,
		ProcessUser:   o.ProcessUser,
		Verbosity:     helpers.CopyIntPtr(o.Verbosity),
		Flags:         helpers.CopyStringSlice(o.Flags),
	}
}

// mergeWith merges the other settings into any
// unset field of the receiver settings object.
func (o *OpenVPN) mergeWith(other OpenVPN) {
	o.Version = helpers.MergeWithString(o.Version, other.Version)
	o.User = helpers.MergeWithStringPtr(o.User, other.User)
	o.Password = helpers.MergeWithStringPtr(o.Password, other.Password)
	o.ConfFile = helpers.MergeWithStringPtr(o.ConfFile, other.ConfFile)
	o.Ciphers = helpers.MergeStringSlices(o.Ciphers, other.Ciphers)
	o.Auth = helpers.MergeWithStringPtr(o.Auth, other.Auth)
	o.Cert = helpers.MergeWithStringPtr(o.Cert, other.Cert)
	o.Key = helpers.MergeWithStringPtr(o.Key, other.Key)
	o.EncryptedKey = helpers.MergeWithStringPtr(o.EncryptedKey, other.EncryptedKey)
	o.KeyPassphrase = helpers.MergeWithStringPtr(o.KeyPassphrase, other.KeyPassphrase)
	o.PIAEncPreset = helpers.MergeWithStringPtr(o.PIAEncPreset, other.PIAEncPreset)
	o.MSSFix = helpers.MergeWithUint16(o.MSSFix, other.MSSFix)
	o.Interface = helpers.MergeWithString(o.Interface, other.Interface)
	o.ProcessUser = helpers.MergeWithString(o.ProcessUser, other.ProcessUser)
	o.Verbosity = helpers.MergeWithIntPtr(o.Verbosity, other.Verbosity)
	o.Flags = helpers.MergeStringSlices(o.Flags, other.Flags)
}

// overrideWith overrides fields of the receiver
// settings object with any field set in the other
// settings.
func (o *OpenVPN) overrideWith(other OpenVPN) {
	o.Version = helpers.OverrideWithString(o.Version, other.Version)
	o.User = helpers.OverrideWithStringPtr(o.User, other.User)
	o.Password = helpers.OverrideWithStringPtr(o.Password, other.Password)
	o.ConfFile = helpers.OverrideWithStringPtr(o.ConfFile, other.ConfFile)
	o.Ciphers = helpers.OverrideWithStringSlice(o.Ciphers, other.Ciphers)
	o.Auth = helpers.OverrideWithStringPtr(o.Auth, other.Auth)
	o.Cert = helpers.OverrideWithStringPtr(o.Cert, other.Cert)
	o.Key = helpers.OverrideWithStringPtr(o.Key, other.Key)
	o.EncryptedKey = helpers.OverrideWithStringPtr(o.EncryptedKey, other.EncryptedKey)
	o.KeyPassphrase = helpers.OverrideWithStringPtr(o.KeyPassphrase, other.KeyPassphrase)
	o.PIAEncPreset = helpers.OverrideWithStringPtr(o.PIAEncPreset, other.PIAEncPreset)
	o.MSSFix = helpers.OverrideWithUint16(o.MSSFix, other.MSSFix)
	o.Interface = helpers.OverrideWithString(o.Interface, other.Interface)
	o.ProcessUser = helpers.OverrideWithString(o.ProcessUser, other.ProcessUser)
	o.Verbosity = helpers.OverrideWithIntPtr(o.Verbosity, other.Verbosity)
	o.Flags = helpers.OverrideWithStringSlice(o.Flags, other.Flags)
}

func (o *OpenVPN) setDefaults(vpnProvider string) {
	o.Version = helpers.DefaultString(o.Version, openvpn.Openvpn25)
	o.User = helpers.DefaultStringPtr(o.User, "")
	if vpnProvider == providers.Mullvad {
		o.Password = helpers.DefaultStringPtr(o.Password, "m")
	} else {
		o.Password = helpers.DefaultStringPtr(o.Password, "")
	}

	o.ConfFile = helpers.DefaultStringPtr(o.ConfFile, "")
	o.Auth = helpers.DefaultStringPtr(o.Auth, "")
	o.Cert = helpers.DefaultStringPtr(o.Cert, "")
	o.Key = helpers.DefaultStringPtr(o.Key, "")
	o.EncryptedKey = helpers.DefaultStringPtr(o.EncryptedKey, "")
	o.KeyPassphrase = helpers.DefaultStringPtr(o.KeyPassphrase, "")

	var defaultEncPreset string
	if vpnProvider == providers.PrivateInternetAccess {
		defaultEncPreset = presets.Strong
	}
	o.PIAEncPreset = helpers.DefaultStringPtr(o.PIAEncPreset, defaultEncPreset)
	o.MSSFix = helpers.DefaultUint16(o.MSSFix, 0)
	o.Interface = helpers.DefaultString(o.Interface, "tun0")
	o.ProcessUser = helpers.DefaultString(o.ProcessUser, "root")
	o.Verbosity = helpers.DefaultInt(o.Verbosity, 1)
}

func (o OpenVPN) String() string {
	return o.toLinesNode().String()
}

func (o OpenVPN) toLinesNode() (node *gotree.Node) {
	node = gotree.New("OpenVPN settings:")
	node.Appendf("OpenVPN version: %s", o.Version)
	node.Appendf("User: %s", helpers.ObfuscatePassword(*o.User))
	node.Appendf("Password: %s", helpers.ObfuscatePassword(*o.Password))

	if *o.ConfFile != "" {
		node.Appendf("Custom configuration file: %s", *o.ConfFile)
	}

	if len(o.Ciphers) > 0 {
		node.Appendf("Ciphers: %s", o.Ciphers)
	}

	if *o.Auth != "" {
		node.Appendf("Auth: %s", *o.Auth)
	}

	if *o.Cert != "" {
		node.Appendf("Client crt: %s", helpers.ObfuscateData(*o.Cert))
	}

	if *o.Key != "" {
		node.Appendf("Client key: %s", helpers.ObfuscateData(*o.Key))
	}

	if *o.EncryptedKey != "" {
		node.Appendf("Encrypted key: %s (key passhrapse %s)",
			helpers.ObfuscateData(*o.EncryptedKey), helpers.ObfuscatePassword(*o.KeyPassphrase))
	}

	if *o.PIAEncPreset != "" {
		node.Appendf("Private Internet Access encryption preset: %s", *o.PIAEncPreset)
	}

	if *o.MSSFix > 0 {
		node.Appendf("MSS Fix: %d", *o.MSSFix)
	}

	if o.Interface != "" {
		node.Appendf("Network interface: %s", o.Interface)
	}

	node.Appendf("Run OpenVPN as: %s", o.ProcessUser)

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
