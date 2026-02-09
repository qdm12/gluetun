package settings

import (
	"fmt"
	"os"

	"github.com/qdm12/gosettings"
	"github.com/qdm12/gosettings/reader"
	"github.com/qdm12/gosettings/validate"
	"github.com/qdm12/gotree"
)

// Socks5 contains settings to configure the SOCKS5 server.
type Socks5 struct {
	// Enabled is true if the server should be running.
	// It defaults to false, and cannot be nil in the internal state.
	Enabled *bool
	ListeningAddress string
	User *string
	Password *string
	Log *bool
}

func (s Socks5) validate() (err error) {
	err = validate.ListeningAddress(s.ListeningAddress, os.Getuid())
	if err != nil {
		return fmt.Errorf("%w: %s", ErrServerAddressNotValid, s.ListeningAddress)
	}
	return nil
}

func (s *Socks5) copy() (copied Socks5) {
	return Socks5{
		Enabled:          gosettings.CopyPointer(s.Enabled),
		ListeningAddress: s.ListeningAddress,
		User:             gosettings.CopyPointer(s.User),
		Password:         gosettings.CopyPointer(s.Password),
		Log:              gosettings.CopyPointer(s.Log),
	}
}

// overrideWith overrides fields of the receiver
// settings object with any field set in the other
// settings.
func (s *Socks5) overrideWith(other Socks5) {
	s.Enabled = gosettings.OverrideWithPointer(s.Enabled, other.Enabled)
	s.ListeningAddress = gosettings.OverrideWithComparable(s.ListeningAddress, other.ListeningAddress)
	s.User = gosettings.OverrideWithPointer(s.User, other.User)
	s.Password = gosettings.OverrideWithPointer(s.Password, other.Password)
	s.Log = gosettings.OverrideWithPointer(s.Log, other.Log)
}

func (s *Socks5) setDefaults() {
	s.Enabled = gosettings.DefaultPointer(s.Enabled, false)
	s.ListeningAddress = gosettings.DefaultComparable(s.ListeningAddress, ":1080")
	s.User = gosettings.DefaultPointer(s.User, "")
	s.Password = gosettings.DefaultPointer(s.Password, "")
	s.Log = gosettings.DefaultPointer(s.Log, false)
}

func (s Socks5) String() string {
	return s.toLinesNode().String()
}

func (s Socks5) toLinesNode() (node *gotree.Node) {
	node = gotree.New("SOCKS5 server settings:")

	node.Appendf("Enabled: %s", gosettings.BoolToYesNo(s.Enabled))
	if !*s.Enabled {
		return node
	}

	node.Appendf("Listening address: %s", s.ListeningAddress)
	node.Appendf("User: %s", *s.User)
	node.Appendf("Password: %s", gosettings.ObfuscateKey(*s.Password))
	node.Appendf("Log: %s", gosettings.BoolToYesNo(s.Log))

	return node
}

func (s *Socks5) read(r *reader.Reader) (err error) {
	s.Enabled, err = r.BoolPtr("SOCKS5")
	if err != nil {
		return err
	}

	s.ListeningAddress = r.String("SOCKS5_LISTENING_ADDRESS")

	s.User = r.Get("SOCKS5_USER",
		reader.ForceLowercase(false))
	s.Password = r.Get("SOCKS5_PASSWORD",
		reader.ForceLowercase(false))

	s.Log, err = r.BoolPtr("SOCKS5_LOG")
	if err != nil {
		return err
	}

	return nil
}
