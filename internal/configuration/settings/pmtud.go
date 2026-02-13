package settings

import (
	"errors"
	"fmt"
	"net/netip"
	"strings"

	"github.com/qdm12/gosettings"
	"github.com/qdm12/gosettings/reader"
	"github.com/qdm12/gotree"
)

// PMTUD contains settings to configure Path MTU Discovery.
type PMTUD struct {
	// Addresses is the redundancy list of addresses to use for path MTU discovery.
	// Each address should either have:
	// - a listening TCP server listening on ip:port
	// - handle the ICMP protocol on ip
	// It cannot be nil in the internal state.
	Addresses []netip.AddrPort `json:"tcp_addresses"`
}

var ErrPMTUDAddressNotValid = errors.New("PMTUD address is not valid")

// Validate validates PMTUD settings.
func (p PMTUD) validate() (err error) {
	for i, addr := range p.Addresses {
		if !addr.IsValid() {
			return fmt.Errorf("%w: at index %d", ErrPMTUDAddressNotValid, i)
		}
	}
	return nil
}

func (p *PMTUD) copy() (copied PMTUD) {
	return PMTUD{
		Addresses: gosettings.CopySlice(p.Addresses),
	}
}

func (p *PMTUD) overrideWith(other PMTUD) {
	p.Addresses = gosettings.OverrideWithSlice(p.Addresses, other.Addresses)
}

func (p *PMTUD) setDefaults() {
	const tlsPort = 443
	defaultTCPAddresses := []netip.AddrPort{
		netip.AddrPortFrom(netip.AddrFrom4([4]byte{1, 1, 1, 1}), tlsPort),
		netip.AddrPortFrom(netip.AddrFrom4([4]byte{8, 8, 8, 8}), tlsPort),
	}
	p.Addresses = gosettings.DefaultSlice(p.Addresses, defaultTCPAddresses)
}

func (p PMTUD) String() string {
	return p.toLinesNode().String()
}

func (p PMTUD) toLinesNode() (node *gotree.Node) {
	node = gotree.New("Path MTU discovery:")
	addrs := make([]string, len(p.Addresses))
	for i, addr := range p.Addresses {
		addrs[i] = addr.String()
	}
	node.Appendf("Addresses: %s", strings.Join(addrs, ", "))
	return node
}

func (p *PMTUD) read(r *reader.Reader) (err error) {
	p.Addresses, err = r.CSVNetipAddrPorts("PMTUD_ADDRESSES")
	if err != nil {
		return err
	}
	return nil
}
