package settings

import (
	"fmt"
	"os"
	"time"

	"github.com/qdm12/gluetun/internal/configuration/settings/helpers"
	"github.com/qdm12/gotree"
	"github.com/qdm12/govalid/address"
)

// Health contains settings for the healthcheck and health server.
type Health struct {
	// ServerAddress is the listening address
	// for the health check server.
	// It cannot be the empty string in the internal state.
	ServerAddress string
	// ReadHeaderTimeout is the HTTP server header read timeout
	// duration of the HTTP server. It defaults to 100 milliseconds.
	ReadHeaderTimeout time.Duration
	// ReadTimeout is the HTTP read timeout duration of the
	//  HTTP server. It defaults to 500 milliseconds.
	ReadTimeout time.Duration
	// TargetAddress is the address (host or host:port)
	// to TCP dial to periodically for the health check.
	// It cannot be the empty string in the internal state.
	TargetAddress string
	VPN           HealthyWait
}

func (h Health) Validate() (err error) {
	uid := os.Getuid()
	_, err = address.Validate(h.ServerAddress,
		address.OptionListening(uid))
	if err != nil {
		return fmt.Errorf("server listening address is not valid: %w", err)
	}

	err = h.VPN.validate()
	if err != nil {
		return fmt.Errorf("health VPN settings: %w", err)
	}

	return nil
}

func (h *Health) copy() (copied Health) {
	return Health{
		ServerAddress:     h.ServerAddress,
		ReadHeaderTimeout: h.ReadHeaderTimeout,
		ReadTimeout:       h.ReadTimeout,
		TargetAddress:     h.TargetAddress,
		VPN:               h.VPN.copy(),
	}
}

// MergeWith merges the other settings into any
// unset field of the receiver settings object.
func (h *Health) MergeWith(other Health) {
	h.ServerAddress = helpers.MergeWithString(h.ServerAddress, other.ServerAddress)
	h.ReadHeaderTimeout = helpers.MergeWithDuration(h.ReadHeaderTimeout, other.ReadHeaderTimeout)
	h.ReadTimeout = helpers.MergeWithDuration(h.ReadTimeout, other.ReadTimeout)
	h.TargetAddress = helpers.MergeWithString(h.TargetAddress, other.TargetAddress)
	h.VPN.mergeWith(other.VPN)
}

// OverrideWith overrides fields of the receiver
// settings object with any field set in the other
// settings.
func (h *Health) OverrideWith(other Health) {
	h.ServerAddress = helpers.OverrideWithString(h.ServerAddress, other.ServerAddress)
	h.ReadHeaderTimeout = helpers.OverrideWithDuration(h.ReadHeaderTimeout, other.ReadHeaderTimeout)
	h.ReadTimeout = helpers.OverrideWithDuration(h.ReadTimeout, other.ReadTimeout)
	h.TargetAddress = helpers.OverrideWithString(h.TargetAddress, other.TargetAddress)
	h.VPN.overrideWith(other.VPN)
}

func (h *Health) SetDefaults() {
	h.ServerAddress = helpers.DefaultString(h.ServerAddress, "127.0.0.1:9999")
	const defaultReadHeaderTimeout = 100 * time.Millisecond
	h.ReadHeaderTimeout = helpers.DefaultDuration(h.ReadHeaderTimeout, defaultReadHeaderTimeout)
	const defaultReadTimeout = 500 * time.Millisecond
	h.ReadTimeout = helpers.DefaultDuration(h.ReadTimeout, defaultReadTimeout)
	h.TargetAddress = helpers.DefaultString(h.TargetAddress, "cloudflare.com:443")
	h.VPN.setDefaults()
}

func (h Health) String() string {
	return h.toLinesNode().String()
}

func (h Health) toLinesNode() (node *gotree.Node) {
	node = gotree.New("Health settings:")
	node.Appendf("Server listening address: %s", h.ServerAddress)
	node.Appendf("Target address: %s", h.TargetAddress)
	node.Appendf("Read header timeout: %s", h.ReadHeaderTimeout)
	node.Appendf("Read timeout: %s", h.ReadTimeout)
	node.AppendNode(h.VPN.toLinesNode("VPN"))
	return node
}
