package subnet

import (
	"net/netip"
)

func FindSubnetsToChange(oldSubnets, newSubnets []netip.Prefix) (subnetsToAdd, subnetsToRemove []netip.Prefix) {
	subnetsToAdd = findSubnetsToAdd(oldSubnets, newSubnets)
	subnetsToRemove = findSubnetsToRemove(oldSubnets, newSubnets)
	return subnetsToAdd, subnetsToRemove
}

func findSubnetsToAdd(oldSubnets, newSubnets []netip.Prefix) (subnetsToAdd []netip.Prefix) {
	for _, newSubnet := range newSubnets {
		found := false
		for _, oldSubnet := range oldSubnets {
			if oldSubnet.String() == newSubnet.String() {
				found = true
				break
			}
		}
		if !found {
			subnetsToAdd = append(subnetsToAdd, newSubnet)
		}
	}
	return subnetsToAdd
}

func findSubnetsToRemove(oldSubnets, newSubnets []netip.Prefix) (subnetsToRemove []netip.Prefix) {
	for _, oldSubnet := range oldSubnets {
		found := false
		for _, newSubnet := range newSubnets {
			if oldSubnet.String() == newSubnet.String() {
				found = true
				break
			}
		}
		if !found {
			subnetsToRemove = append(subnetsToRemove, oldSubnet)
		}
	}
	return subnetsToRemove
}

func RemoveSubnetFromSubnets(subnets []netip.Prefix, subnet netip.Prefix) []netip.Prefix {
	L := len(subnets)
	for i := range subnets {
		if subnet.String() == subnets[i].String() {
			subnets[i] = subnets[L-1]
			subnets = subnets[:L-1]
			break
		}
	}
	return subnets
}
