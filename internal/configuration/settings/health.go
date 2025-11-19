package settings

import (
	"fmt"
	"net/netip"
	"os"

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
	// TargetAddresses are the addresses (host or host:port)
	// to TCP TLS dial to periodically for the health check.
	// Addresses after the first one are used as fallbacks for retries.
	// It cannot be empty in the internal state.
	TargetAddresses []string
	// ICMPTargetIP is the IP address to use for ICMP echo requests
	// in the health checker. It can be set to an unspecified address (0.0.0.0)
	// such that the VPN server IP is used, although this can be less reliable.
	// It defaults to 1.1.1.1, and cannot be left empty (invalid) in the internal state.
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
		ServerAddress:   h.ServerAddress,
		TargetAddresses: h.TargetAddresses,
		ICMPTargetIP:    h.ICMPTargetIP,
		RestartVPN:      gosettings.CopyPointer(h.RestartVPN),
	}
}

// OverrideWith overrides fields of the receiver
// settings object with any field set in the other
// settings.
func (h *Health) OverrideWith(other Health) {
	h.ServerAddress = gosettings.OverrideWithComparable(h.ServerAddress, other.ServerAddress)
	h.TargetAddresses = gosettings.OverrideWithSlice(h.TargetAddresses, other.TargetAddresses)
	h.ICMPTargetIP = gosettings.OverrideWithComparable(h.ICMPTargetIP, other.ICMPTargetIP)
	h.RestartVPN = gosettings.OverrideWithPointer(h.RestartVPN, other.RestartVPN)
}

func (h *Health) SetDefaults() {
	h.ServerAddress = gosettings.DefaultComparable(h.ServerAddress, "127.0.0.1:9999")
	h.TargetAddresses = gosettings.DefaultSlice(h.TargetAddresses, []string{"cloudflare.com:443", "github.com:443"})
	h.ICMPTargetIP = gosettings.DefaultComparable(h.ICMPTargetIP, netip.AddrFrom4([4]byte{1, 1, 1, 1}))
	h.RestartVPN = gosettings.DefaultPointer(h.RestartVPN, true)
}

func (h Health) String() string {
	return h.toLinesNode().String()
}

func (h Health) toLinesNode() (node *gotree.Node) {
	node = gotree.New("Health settings:")
	node.Appendf("Server listening address: %s", h.ServerAddress)
	targetAddrs := node.Appendf("Target addresses:")
	for _, targetAddr := range h.TargetAddresses {
		targetAddrs.Append(targetAddr)
	}
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
	h.TargetAddresses = r.CSV("HEALTH_TARGET_ADDRESSES",
		reader.RetroKeys("HEALTH_ADDRESS_TO_PING", "HEALTH_TARGET_ADDRESS"))
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
