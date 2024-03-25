package settings

import (
	"fmt"

	"github.com/qdm12/gosettings"
	"github.com/qdm12/gosettings/reader"
	"github.com/qdm12/gotree"
	"github.com/qdm12/ss-server/pkg/tcpudp"
)

// Shadowsocks contains settings to configure the Shadowsocks server.
type Shadowsocks struct {
	// Enabled is true if the server should be running.
	// It defaults to false, and cannot be nil in the internal state.
	Enabled *bool
	// Settings are settings for the TCP+UDP server.
	tcpudp.Settings
}

func (s Shadowsocks) validate() (err error) {
	return s.Settings.Validate()
}

func (s *Shadowsocks) copy() (copied Shadowsocks) {
	return Shadowsocks{
		Enabled:  gosettings.CopyPointer(s.Enabled),
		Settings: s.Settings.Copy(),
	}
}

// overrideWith overrides fields of the receiver
// settings object with any field set in the other
// settings.
func (s *Shadowsocks) overrideWith(other Shadowsocks) {
	s.Enabled = gosettings.OverrideWithPointer(s.Enabled, other.Enabled)
	s.Settings.OverrideWith(other.Settings)
}

func (s *Shadowsocks) setDefaults() {
	s.Enabled = gosettings.DefaultPointer(s.Enabled, false)
	s.Settings.SetDefaults()
}

func (s Shadowsocks) String() string {
	return s.toLinesNode().String()
}

func (s Shadowsocks) toLinesNode() (node *gotree.Node) {
	node = gotree.New("Shadowsocks server settings:")

	node.Appendf("Enabled: %s", gosettings.BoolToYesNo(s.Enabled))
	if !*s.Enabled {
		return node
	}

	// TODO have ToLinesNode in qdm12/ss-server
	node.Appendf("Listening address: %s", *s.Settings.Address)
	node.Appendf("Cipher: %s", s.Settings.CipherName)
	node.Appendf("Password: %s", gosettings.ObfuscateKey(*s.Settings.Password))
	node.Appendf("Log addresses: %s", gosettings.BoolToYesNo(s.Settings.LogAddresses))

	return node
}

func (s *Shadowsocks) read(r *reader.Reader) (err error) {
	s.Enabled, err = r.BoolPtr("SHADOWSOCKS")
	if err != nil {
		return err
	}

	s.Settings.Address, err = readShadowsocksAddress(r)
	if err != nil {
		return err
	}
	s.Settings.LogAddresses, err = r.BoolPtr("SHADOWSOCKS_LOG")
	if err != nil {
		return err
	}
	s.Settings.CipherName = r.String("SHADOWSOCKS_CIPHER",
		reader.RetroKeys("SHADOWSOCKS_METHOD"))
	s.Settings.Password = r.Get("SHADOWSOCKS_PASSWORD",
		reader.ForceLowercase(false))

	return nil
}

func readShadowsocksAddress(r *reader.Reader) (address *string, err error) {
	const currentKey = "SHADOWSOCKS_LISTENING_ADDRESS"
	port, err := r.Uint16Ptr("SHADOWSOCKS_PORT", reader.IsRetro(currentKey)) // retro-compatibility
	if err != nil {
		return nil, err
	} else if port != nil {
		return ptrTo(fmt.Sprintf(":%d", *port)), nil
	}

	return r.Get(currentKey), nil
}
