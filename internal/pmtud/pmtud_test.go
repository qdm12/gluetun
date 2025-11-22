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
			mtus: []int{0},
		},
		"0_1": {
			maxMTU: 1,
			mtus:   []int{0, 1},
		},
		"0_8": {
			maxMTU: 8,
			mtus:   []int{0, 1, 2, 3, 4, 5, 6, 7, 8},
		},
		"0_12": {
			maxMTU: 12,
			mtus:   []int{0, 1, 2, 4, 5, 6, 7, 8, 10, 11, 12},
		},
		"0_80": {
			maxMTU: 80,
			mtus:   []int{0, 8, 16, 24, 32, 40, 48, 56, 64, 72, 80},
		},
		"0_100": {
			maxMTU: 100,
			mtus:   []int{0, 10, 20, 30, 40, 50, 60, 70, 80, 90, 100},
		},
		"1280_1500": {
			minMTU: 1280,
			maxMTU: 1500,
			mtus:   []int{1280, 1302, 1324, 1346, 1368, 1390, 1412, 1434, 1456, 1478, 1500},
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
