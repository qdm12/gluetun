package ip

import (
	"net/netip"
	"slices"

	"github.com/qdm12/gluetun/internal/pmtud/constants"
)

func GetFamilies(dsts []netip.AddrPort) (families []int) {
	const maxFamilies = 2
	families = make([]int, 0, maxFamilies)
	for _, dst := range dsts {
		family := GetFamily(dst)
		if !slices.Contains(families, family) {
			families = append(families, family)
		}
	}
	return families
}

func GetFamily(dst netip.AddrPort) int {
	if dst.Addr().Is4() {
		return constants.AF_INET
	}
	return constants.AF_INET6
}
