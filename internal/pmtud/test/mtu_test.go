package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_MakeMTUsToTest(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		minMTU uint32
		maxMTU uint32
		mtus   []uint32
	}{
		"0_0": {
			mtus: []uint32{0},
		},
		"0_1": {
			maxMTU: 1,
			mtus:   []uint32{0, 1},
		},
		"0_8": {
			maxMTU: 8,
			mtus:   []uint32{0, 1, 2, 3, 4, 5, 6, 7, 8},
		},
		"0_12": {
			maxMTU: 12,
			mtus:   []uint32{0, 1, 2, 4, 5, 6, 7, 8, 10, 11, 12},
		},
		"0_80": {
			maxMTU: 80,
			mtus:   []uint32{0, 8, 16, 24, 32, 40, 48, 56, 64, 72, 80},
		},
		"0_100": {
			maxMTU: 100,
			mtus:   []uint32{0, 10, 20, 30, 40, 50, 60, 70, 80, 90, 100},
		},
		"1280_1500": {
			minMTU: 1280,
			maxMTU: 1500,
			mtus:   []uint32{1280, 1302, 1324, 1346, 1368, 1390, 1412, 1434, 1456, 1478, 1500},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			mtus := MakeMTUsToTest(testCase.minMTU, testCase.maxMTU)
			assert.Equal(t, testCase.mtus, mtus)
		})
	}
}
