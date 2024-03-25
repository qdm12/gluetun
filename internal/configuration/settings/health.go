package settings

import (
	"fmt"
	"os"
	"time"

	"github.com/qdm12/gosettings"
	"github.com/qdm12/gosettings/reader"
	"github.com/qdm12/gosettings/validate"
	"github.com/qdm12/gotree"
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
	// HTTP server. It defaults to 500 milliseconds.
	ReadTimeout time.Duration
	// TargetAddress is the address (host or host:port)
	// to TCP dial to periodically for the health check.
	// It cannot be the empty string in the internal state.
	TargetAddress string
	// SuccessWait is the duration to wait to re-run the
	// healthcheck after a successful healthcheck.
	// It defaults to 5 seconds and cannot be zero in
	// the internal state.
	SuccessWait time.Duration
	// VPN has health settings specific to the VPN loop.
	VPN HealthyWait
}

func (h Health) Validate() (err error) {
	err = validate.ListeningAddress(h.ServerAddress, os.Getuid())
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
		SuccessWait:       h.SuccessWait,
		VPN:               h.VPN.copy(),
	}
}

// OverrideWith overrides fields of the receiver
// settings object with any field set in the other
// settings.
func (h *Health) OverrideWith(other Health) {
	h.ServerAddress = gosettings.OverrideWithComparable(h.ServerAddress, other.ServerAddress)
	h.ReadHeaderTimeout = gosettings.OverrideWithComparable(h.ReadHeaderTimeout, other.ReadHeaderTimeout)
	h.ReadTimeout = gosettings.OverrideWithComparable(h.ReadTimeout, other.ReadTimeout)
	h.TargetAddress = gosettings.OverrideWithComparable(h.TargetAddress, other.TargetAddress)
	h.SuccessWait = gosettings.OverrideWithComparable(h.SuccessWait, other.SuccessWait)
	h.VPN.overrideWith(other.VPN)
}

func (h *Health) SetDefaults() {
	h.ServerAddress = gosettings.DefaultComparable(h.ServerAddress, "127.0.0.1:9999")
	const defaultReadHeaderTimeout = 100 * time.Millisecond
	h.ReadHeaderTimeout = gosettings.DefaultComparable(h.ReadHeaderTimeout, defaultReadHeaderTimeout)
	const defaultReadTimeout = 500 * time.Millisecond
	h.ReadTimeout = gosettings.DefaultComparable(h.ReadTimeout, defaultReadTimeout)
	h.TargetAddress = gosettings.DefaultComparable(h.TargetAddress, "cloudflare.com:443")
	const defaultSuccessWait = 5 * time.Second
	h.SuccessWait = gosettings.DefaultComparable(h.SuccessWait, defaultSuccessWait)
	h.VPN.setDefaults()
}

func (h Health) String() string {
	return h.toLinesNode().String()
}

func (h Health) toLinesNode() (node *gotree.Node) {
	node = gotree.New("Health settings:")
	node.Appendf("Server listening address: %s", h.ServerAddress)
	node.Appendf("Target address: %s", h.TargetAddress)
	node.Appendf("Duration to wait after success: %s", h.SuccessWait)
	node.Appendf("Read header timeout: %s", h.ReadHeaderTimeout)
	node.Appendf("Read timeout: %s", h.ReadTimeout)
	node.AppendNode(h.VPN.toLinesNode("VPN"))
	return node
}

func (h *Health) Read(r *reader.Reader) (err error) {
	h.ServerAddress = r.String("HEALTH_SERVER_ADDRESS")
	h.TargetAddress = r.String("HEALTH_TARGET_ADDRESS",
		reader.RetroKeys("HEALTH_ADDRESS_TO_PING"))

	h.SuccessWait, err = r.Duration("HEALTH_SUCCESS_WAIT_DURATION")
	if err != nil {
		return err
	}

	err = h.VPN.read(r)
	if err != nil {
		return fmt.Errorf("VPN health settings: %w", err)
	}

	return nil
}
