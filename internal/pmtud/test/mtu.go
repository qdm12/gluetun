package test

import (
	"fmt"
	"math"
)

// MakeMTUsToTest determines a slice of MTU values to test
// between minMTU and maxMTU inclusive. It creates an MTU
// slice of length up to 11 MTUs such that:
// - the first element is the minMTU
// - the last element is the maxMTU
// - elements in-between are separated as close to each other
// The number 11 is chosen to find the final MTU in 3 searches,
// with a total search space of 1728 MTUs which is enough;
// to find it in 2 searches requires 37 parallel queries which
// could be blocked by firewalls.
func MakeMTUsToTest(minMTU, maxMTU uint32) (mtus []uint32) {
	const mtusLength = 11 // find the final MTU in 3 searches
	diff := maxMTU - minMTU
	switch {
	case minMTU > maxMTU:
		panic(fmt.Sprintf("minMTU %d is greater than maxMTU %d", minMTU, maxMTU))
	case diff <= mtusLength:
		mtus = make([]uint32, 0, diff)
		for mtu := minMTU; mtu <= maxMTU; mtu++ {
			mtus = append(mtus, mtu)
		}
	default:
		step := float64(diff) / float64(mtusLength-1)
		mtus = make([]uint32, 0, mtusLength)
		for mtu := float64(minMTU); len(mtus) < mtusLength-1; mtu += step {
			mtus = append(mtus, uint32(math.Round(mtu)))
		}
		mtus = append(mtus, maxMTU) // last element is the maxMTU
	}

	return mtus
}
