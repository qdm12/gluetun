package settings

import (
	"fmt"
	"net/netip"
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
	// to TCP TLS dial to periodically for the health check.
	// It cannot be the empty string in the internal state.
	TargetAddress string
	// ICMPTargetIP is the IP address to use for ICMP echo requests
	// in the health checker. It can be set to an unspecified address (0.0.0.0)
	// such that the VPN server IP is used, which is also the default behavior.
	ICMPTargetIP netip.Addr
	// RestartVPN indicates whether to restart the VPN connection
	// when the healthcheck fails.
	RestartVPN *bool
}

func (h Health) Validate() (err error) {
	err = validate.ListeningAddress(h.ServerAddress, os.Getuid())
	if err != nil {
		return fmt.Errorf("server listening address is not valid: %w", err)
	}

	return nil
}

func (h *Health) copy() (copied Health) {
	return Health{
		ServerAddress:     h.ServerAddress,
		ReadHeaderTimeout: h.ReadHeaderTimeout,
		ReadTimeout:       h.ReadTimeout,
		TargetAddress:     h.TargetAddress,
		ICMPTargetIP:      h.ICMPTargetIP,
		RestartVPN:        gosettings.CopyPointer(h.RestartVPN),
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
	h.ICMPTargetIP = gosettings.OverrideWithComparable(h.ICMPTargetIP, other.ICMPTargetIP)
	h.RestartVPN = gosettings.OverrideWithPointer(h.RestartVPN, other.RestartVPN)
}

func (h *Health) SetDefaults() {
	h.ServerAddress = gosettings.DefaultComparable(h.ServerAddress, "127.0.0.1:9999")
	const defaultReadHeaderTimeout = 100 * time.Millisecond
	h.ReadHeaderTimeout = gosettings.DefaultComparable(h.ReadHeaderTimeout, defaultReadHeaderTimeout)
	const defaultReadTimeout = 500 * time.Millisecond
	h.ReadTimeout = gosettings.DefaultComparable(h.ReadTimeout, defaultReadTimeout)
	h.TargetAddress = gosettings.DefaultComparable(h.TargetAddress, "cloudflare.com:443")
	h.ICMPTargetIP = gosettings.DefaultComparable(h.ICMPTargetIP, netip.IPv4Unspecified()) // use the VPN server IP
	h.RestartVPN = gosettings.DefaultPointer(h.RestartVPN, true)
}

func (h Health) String() string {
	return h.toLinesNode().String()
}

func (h Health) toLinesNode() (node *gotree.Node) {
	node = gotree.New("Health settings:")
	node.Appendf("Server listening address: %s", h.ServerAddress)
	node.Appendf("Target address: %s", h.TargetAddress)
	icmpTarget := "VPN server IP"
	if !h.ICMPTargetIP.IsUnspecified() {
		icmpTarget = h.ICMPTargetIP.String()
	}
	node.Appendf("ICMP target IP: %s", icmpTarget)
	node.Appendf("Restart VPN on healthcheck failure: %s", gosettings.BoolToYesNo(h.RestartVPN))
	return node
}

func (h *Health) Read(r *reader.Reader) (err error) {
	h.ServerAddress = r.String("HEALTH_SERVER_ADDRESS")
	h.TargetAddress = r.String("HEALTH_TARGET_ADDRESS",
		reader.RetroKeys("HEALTH_ADDRESS_TO_PING"))
	h.ICMPTargetIP, err = r.NetipAddr("HEALTH_ICMP_TARGET_IP")
	if err != nil {
		return err
	}
	h.RestartVPN, err = r.BoolPtr("HEALTH_RESTART_VPN")
	if err != nil {
		return err
	}
	return nil
}
