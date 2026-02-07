package settings

import (
	"errors"
	"fmt"
	"net/netip"

	"github.com/qdm12/gosettings"
	"github.com/qdm12/gosettings/reader"
	"github.com/qdm12/gotree"
)

// PMTUD contains settings to configure Path MTU Discovery.
type PMTUD struct {
	Enabled *bool `json:"enabled"`
	// TCPAddresses is the list of TCP addresses to use for path MTU discovery.
	// It cannot be nil in the internal state.
	TCPAddresses []netip.AddrPort `json:"tcp_addresses"`
}

var ErrPMTUDTCPAddressNotValid = errors.New("PMTUD TCP address is not valid")

// Validate validates PMTUD settings.
func (p PMTUD) validate() (err error) {
	for i, addr := range p.TCPAddresses {
		if !addr.IsValid() {
			return fmt.Errorf("%w: at index %d",
				ErrPMTUDTCPAddressNotValid, i)
		}
	}
	return nil
}

func (p *PMTUD) copy() (copied PMTUD) {
	return PMTUD{
		Enabled:      gosettings.CopyPointer(p.Enabled),
		TCPAddresses: gosettings.CopySlice(p.TCPAddresses),
	}
}

func (p *PMTUD) overrideWith(other PMTUD) {
	p.Enabled = gosettings.OverrideWithPointer(p.Enabled, other.Enabled)
	p.TCPAddresses = gosettings.OverrideWithSlice(p.TCPAddresses, other.TCPAddresses)
}

func (p *PMTUD) setDefaults() {
	const defaultEnabled = true
	p.Enabled = gosettings.DefaultPointer(p.Enabled, defaultEnabled)
	const tlsPort = 443
	defaultTCPAddresses := []netip.AddrPort{
		netip.AddrPortFrom(netip.AddrFrom4([4]byte{1, 1, 1, 1}), tlsPort),
		netip.AddrPortFrom(netip.AddrFrom4([4]byte{8, 8, 8, 8}), tlsPort),
	}
	p.TCPAddresses = gosettings.DefaultSlice(p.TCPAddresses, defaultTCPAddresses)
}

func (p PMTUD) String() string {
	return p.toLinesNode().String()
}

func (p PMTUD) toLinesNode() (node *gotree.Node) {
	if !*p.Enabled {
		return gotree.New("Path MTU discovery: disabled")
	}
	node = gotree.New("Path MTU discovery:")
	tcpNode := node.Append("TCP addresses:")
	for _, addr := range p.TCPAddresses {
		tcpNode.Append(addr.String())
	}
	return node
}

func (p *PMTUD) read(r *reader.Reader) (err error) {
	p.Enabled, err = r.BoolPtr("PMTUD")
	if err != nil {
		return err
	}

	p.TCPAddresses, err = r.CSVNetipAddrPorts("PMTUD_TCP_ADDRESSES")
	if err != nil {
		return err
	}
	return nil
}
