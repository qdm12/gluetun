package settings

import (
	"fmt"
	"os"

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
		return fmt.Errorf("%w: %s",
			ErrServerAddressNotValid, err)
	}

	err = h.VPN.validate()
	if err != nil {
		return fmt.Errorf("health VPN settings validation failed: %w", err)
	}

	return nil
}

func (h *Health) copy() (copied Health) {
	return Health{
		ServerAddress: h.ServerAddress,
		TargetAddress: h.TargetAddress,
		VPN:           h.VPN.copy(),
	}
}

// MergeWith merges the other settings into any
// unset field of the receiver settings object.
func (h *Health) MergeWith(other Health) {
	h.ServerAddress = helpers.MergeWithString(h.ServerAddress, other.ServerAddress)
	h.TargetAddress = helpers.MergeWithString(h.TargetAddress, other.TargetAddress)
	h.VPN.mergeWith(other.VPN)
}

// OverrideWith overrides fields of the receiver
// settings object with any field set in the other
// settings.
func (h *Health) OverrideWith(other Health) {
	h.ServerAddress = helpers.OverrideWithString(h.ServerAddress, other.ServerAddress)
	h.TargetAddress = helpers.OverrideWithString(h.TargetAddress, other.TargetAddress)
	h.VPN.overrideWith(other.VPN)
}

func (h *Health) SetDefaults() {
	h.ServerAddress = helpers.DefaultString(h.ServerAddress, "127.0.0.1:9999")
	h.TargetAddress = helpers.DefaultString(h.TargetAddress, "github.com:443")
	h.VPN.setDefaults()
}

func (h Health) String() string {
	return h.toLinesNode().String()
}

func (h Health) toLinesNode() (node *gotree.Node) {
	node = gotree.New("Health settings:")
	node.Appendf("Server listening address: %s", h.ServerAddress)
	node.Appendf("Target address: %s", h.TargetAddress)
	node.AppendNode(h.VPN.toLinesNode("VPN"))
	return node
}
