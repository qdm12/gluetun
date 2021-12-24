package settings

import (
	"fmt"
	"os"

	"github.com/qdm12/gluetun/internal/configuration/settings/helpers"
	"github.com/qdm12/govalid/address"
)

// Health contains settings for the healthcheck and health server.
type Health struct {
	// ServerAddress is the listening address
	// for the health check server.
	// It cannot be the empty string in the internal state.
	ServerAddress string `json:"server_address,omitempty"`
	// AddressToPing is the IP address or domain name to
	// ping periodically for the health check.
	// It cannot be the empty string in the internal state.
	AddressToPing string      `json:"address_to_ping,omitempty"`
	VPN           HealthyWait `json:"vpn_wait"`
}

func (h Health) validate() (err error) {
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
		AddressToPing: h.AddressToPing,
		VPN:           h.VPN.copy(),
	}
}

// mergeWith merges the other settings into any
// unset field of the receiver settings object.
func (h *Health) mergeWith(other Health) {
	h.ServerAddress = helpers.MergeWithString(h.ServerAddress, other.ServerAddress)
	h.AddressToPing = helpers.MergeWithString(h.AddressToPing, other.AddressToPing)
	h.VPN.mergeWith(other.VPN)
}

// overrideWith overrides fields of the receiver
// settings object with any field set in the other
// settings.
func (h *Health) overrideWith(other Health) {
	h.ServerAddress = helpers.OverrideWithString(h.ServerAddress, other.ServerAddress)
	h.AddressToPing = helpers.OverrideWithString(h.AddressToPing, other.AddressToPing)
	h.VPN.overrideWith(other.VPN)
}

func (h *Health) setDefaults() {
	h.ServerAddress = helpers.DefaultString(h.ServerAddress, "127.0.0.1:9999")
	h.AddressToPing = helpers.DefaultString(h.AddressToPing, "github.com")
	h.VPN.setDefaults()
}
