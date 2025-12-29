package settings

import (
	"net/netip"

	"github.com/qdm12/gosettings"
	"github.com/qdm12/gosettings/reader"
	"github.com/qdm12/gotree"
)

// IPv6 contains settings regarding IPv6 configuration.
type IPv6 struct {
	// CheckAddress is the TCP ip:port address to dial to check
	// IPv6 is supported, in case a default IPv6 route is found.
	// It defaults to cloudflare.com address [2606:4700::6810:84e5]:443
	CheckAddress netip.AddrPort
}

func (i IPv6) validate() (err error) {
	return nil
}

func (i *IPv6) copy() (copied IPv6) {
	return IPv6{
		CheckAddress: i.CheckAddress,
	}
}

func (i *IPv6) overrideWith(other IPv6) {
	i.CheckAddress = gosettings.OverrideWithValidator(i.CheckAddress, other.CheckAddress)
}

func (i *IPv6) setDefaults() {
	defaultCheckAddress := netip.MustParseAddrPort("[2606:4700::6810:84e5]:443")
	i.CheckAddress = gosettings.DefaultComparable(i.CheckAddress, defaultCheckAddress)
}

func (i IPv6) String() string {
	return i.toLinesNode().String()
}

func (i IPv6) toLinesNode() (node *gotree.Node) {
	node = gotree.New("IPv6 settings:")
	node.Appendf("Check address: %s", i.CheckAddress)
	return node
}

func (i *IPv6) read(r *reader.Reader) (err error) {
	i.CheckAddress, err = r.NetipAddrPort("IPV6_CHECK_ADDRESS")
	return err
}
