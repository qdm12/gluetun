package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_portsToString(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		ports []uint16
		s     string
	}{
		"no_port": {
			s: "no port forwarded",
		},
		"one_port": {
			ports: []uint16{123},
			s:     "port forwarded is 123",
		},
		"two_ports": {
			ports: []uint16{123, 456},
			s:     "ports forwarded are 123 and 456",
		},
		"three_ports": {
			ports: []uint16{123, 456, 789},
			s:     "ports forwarded are 123, 456 and 789",
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := portsToString(testCase.ports)

			assert.Equal(t, testCase.s, s)
		})
	}
}
