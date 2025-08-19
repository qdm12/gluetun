package pmtud

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_makeMTUsToTest(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		minMTU int
		maxMTU int
		mtus   []int
	}{
		"0_0": {
			mtus: []int{},
		},
		"0_1": {
			maxMTU: 1,
			mtus:   []int{1},
		},
		"0_8": {
			maxMTU: 8,
			mtus:   []int{1, 2, 3, 4, 5, 6, 7, 8},
		},
		"0_12": {
			maxMTU: 12,
			mtus:   []int{2, 3, 5, 6, 8, 9, 11, 12},
		},
		"0_80": {
			maxMTU: 80,
			mtus:   []int{10, 20, 30, 40, 50, 60, 70, 80},
		},
		"0_100": {
			maxMTU: 100,
			mtus:   []int{12, 24, 36, 48, 60, 72, 84, 100},
		},
		"1280_1500": {
			minMTU: 1280,
			maxMTU: 1500,
			mtus:   []int{1307, 1334, 1361, 1388, 1415, 1442, 1469, 1500},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			mtus := makeMTUsToTest(testCase.minMTU, testCase.maxMTU)
			assert.Equal(t, testCase.mtus, mtus)
		})
	}
}
