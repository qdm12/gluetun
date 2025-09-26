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
			mtus:   []int{1, 2, 3, 4, 5, 7, 8, 9, 10, 11, 12},
		},
		"0_80": {
			maxMTU: 80,
			mtus:   []int{7, 14, 21, 28, 35, 42, 49, 56, 63, 70, 80},
		},
		"0_100": {
			maxMTU: 100,
			mtus:   []int{9, 18, 27, 36, 45, 54, 63, 72, 81, 90, 100},
		},
		"1280_1500": {
			minMTU: 1280,
			maxMTU: 1500,
			mtus:   []int{1300, 1320, 1340, 1360, 1380, 1400, 1420, 1440, 1460, 1480, 1500},
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
