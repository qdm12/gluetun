package settings

import (
	"fmt"
	"os"
	"time"

	"github.com/qdm12/gluetun/internal/configuration/settings/helpers"
	"github.com/qdm12/gotree"
	"github.com/qdm12/govalid/address"
)

// HTTPProxy contains settings to configure the HTTP proxy.
type HTTPProxy struct {
	// User is the username to use for the HTTP proxy.
	// It cannot be nil in the internal state.
	User *string
	// Password is the password to use for the HTTP proxy.
	// It cannot be nil in the internal state.
	Password *string
	// ListeningAddress is the listening address
	// of the HTTP proxy server.
	// It cannot be the empty string in the internal state.
	ListeningAddress string
	// Enabled is true if the HTTP proxy server should run,
	// and false otherwise. It cannot be nil in the
	// internal state.
	Enabled *bool
	// Stealth is true if the HTTP proxy server should hide
	// each request has been proxied to the destination.
	// It cannot be nil in the internal state.
	Stealth *bool
	// Log is true if the HTTP proxy server should log
	// each request/response. It cannot be nil in the
	// internal state.
	Log *bool
	// ReadHeaderTimeout is the HTTP header read timeout duration
	// of the HTTP server. It defaults to 1 second if left unset.
	ReadHeaderTimeout time.Duration
	// ReadTimeout is the HTTP read timeout duration
	// of the HTTP server. It defaults to 3 seconds if left unset.
	ReadTimeout time.Duration
}

func (h HTTPProxy) validate() (err error) {
	// Do not validate user and password

	uid := os.Getuid()
	_, err = address.Validate(h.ListeningAddress, address.OptionListening(uid))
	if err != nil {
		return fmt.Errorf("%w: %s", ErrServerAddressNotValid, h.ListeningAddress)
	}

	return nil
}

func (h *HTTPProxy) copy() (copied HTTPProxy) {
	return HTTPProxy{
		User:              helpers.CopyPointer(h.User),
		Password:          helpers.CopyPointer(h.Password),
		ListeningAddress:  h.ListeningAddress,
		Enabled:           helpers.CopyPointer(h.Enabled),
		Stealth:           helpers.CopyPointer(h.Stealth),
		Log:               helpers.CopyPointer(h.Log),
		ReadHeaderTimeout: h.ReadHeaderTimeout,
		ReadTimeout:       h.ReadTimeout,
	}
}

// mergeWith merges the other settings into any
// unset field of the receiver settings object.
func (h *HTTPProxy) mergeWith(other HTTPProxy) {
	h.User = helpers.MergeWithPointer(h.User, other.User)
	h.Password = helpers.MergeWithPointer(h.Password, other.Password)
	h.ListeningAddress = helpers.MergeWithString(h.ListeningAddress, other.ListeningAddress)
	h.Enabled = helpers.MergeWithPointer(h.Enabled, other.Enabled)
	h.Stealth = helpers.MergeWithPointer(h.Stealth, other.Stealth)
	h.Log = helpers.MergeWithPointer(h.Log, other.Log)
	h.ReadHeaderTimeout = helpers.MergeWithNumber(h.ReadHeaderTimeout, other.ReadHeaderTimeout)
	h.ReadTimeout = helpers.MergeWithNumber(h.ReadTimeout, other.ReadTimeout)
}

// overrideWith overrides fields of the receiver
// settings object with any field set in the other
// settings.
func (h *HTTPProxy) overrideWith(other HTTPProxy) {
	h.User = helpers.OverrideWithPointer(h.User, other.User)
	h.Password = helpers.OverrideWithPointer(h.Password, other.Password)
	h.ListeningAddress = helpers.OverrideWithString(h.ListeningAddress, other.ListeningAddress)
	h.Enabled = helpers.OverrideWithPointer(h.Enabled, other.Enabled)
	h.Stealth = helpers.OverrideWithPointer(h.Stealth, other.Stealth)
	h.Log = helpers.OverrideWithPointer(h.Log, other.Log)
	h.ReadHeaderTimeout = helpers.OverrideWithNumber(h.ReadHeaderTimeout, other.ReadHeaderTimeout)
	h.ReadTimeout = helpers.OverrideWithNumber(h.ReadTimeout, other.ReadTimeout)
}

func (h *HTTPProxy) setDefaults() {
	h.User = helpers.DefaultPointer(h.User, "")
	h.Password = helpers.DefaultPointer(h.Password, "")
	h.ListeningAddress = helpers.DefaultString(h.ListeningAddress, ":8888")
	h.Enabled = helpers.DefaultPointer(h.Enabled, false)
	h.Stealth = helpers.DefaultPointer(h.Stealth, false)
	h.Log = helpers.DefaultPointer(h.Log, false)
	const defaultReadHeaderTimeout = time.Second
	h.ReadHeaderTimeout = helpers.DefaultNumber(h.ReadHeaderTimeout, defaultReadHeaderTimeout)
	const defaultReadTimeout = 3 * time.Second
	h.ReadTimeout = helpers.DefaultNumber(h.ReadTimeout, defaultReadTimeout)
}

func (h HTTPProxy) String() string {
	return h.toLinesNode().String()
}

func (h HTTPProxy) toLinesNode() (node *gotree.Node) {
	node = gotree.New("HTTP proxy settings:")
	node.Appendf("Enabled: %s", helpers.BoolPtrToYesNo(h.Enabled))
	if !*h.Enabled {
		return node
	}

	node.Appendf("Listening address: %s", h.ListeningAddress)
	node.Appendf("User: %s", *h.User)
	node.Appendf("Password: %s", helpers.ObfuscatePassword(*h.Password))
	node.Appendf("Stealth mode: %s", helpers.BoolPtrToYesNo(h.Stealth))
	node.Appendf("Log: %s", helpers.BoolPtrToYesNo(h.Log))
	node.Appendf("Read header timeout: %s", h.ReadHeaderTimeout)
	node.Appendf("Read timeout: %s", h.ReadTimeout)

	return node
}
