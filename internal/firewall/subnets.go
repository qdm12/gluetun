package firewall

import (
	"net"
)

func findSubnetsToAdd(oldSubnets, newSubnets []net.IPNet) (subnetsToAdd []net.IPNet) {
	for _, newSubnet := range newSubnets {
		found := false
		for _, oldSubnet := range oldSubnets {
			if subnetsAreEqual(oldSubnet, newSubnet) {
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

func findSubnetsToRemove(oldSubnets, newSubnets []net.IPNet) (subnetsToRemove []net.IPNet) {
	for _, oldSubnet := range oldSubnets {
		found := false
		for _, newSubnet := range newSubnets {
			if subnetsAreEqual(oldSubnet, newSubnet) {
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

func subnetsAreEqual(a, b net.IPNet) bool {
	return a.IP.Equal(b.IP) && a.Mask.String() == b.Mask.String()
}

func removeSubnetFromSubnets(subnets []net.IPNet, subnet net.IPNet) []net.IPNet {
	L := len(subnets)
	for i := range subnets {
		if subnetsAreEqual(subnet, subnets[i]) {
			subnets[i] = subnets[L-1]
			subnets = subnets[:L-1]
			break
		}
	}
	return subnets
}
