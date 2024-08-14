package firewall

import (
	"net/netip"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseIptablesInstruction(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		s           string
		instruction iptablesInstruction
		errWrapped  error
		errMessage  string
	}{
		"no_instruction": {
			errWrapped: ErrIptablesCommandMalformed,
			errMessage: "iptables command is malformed: empty instruction",
		},
		"uneven_fields": {
			s:          "-A",
			errWrapped: ErrIptablesCommandMalformed,
			errMessage: "iptables command is malformed: fields count 1 is not even: \"-A\"",
		},
		"unknown_key": {
			s:          "-x something",
			errWrapped: ErrIptablesCommandMalformed,
			errMessage: "parsing \"-x something\": iptables command is malformed: unknown key \"-x\"",
		},
		"one_pair": {
			s: "-A INPUT",
			instruction: iptablesInstruction{
				table:  "filter",
				chain:  "INPUT",
				append: true,
			},
		},
		"instruction_A": {
			s: "-A INPUT -i tun0 -p tcp -m tcp -s 1.2.3.4/32 --dport 10000 -j ACCEPT",
			instruction: iptablesInstruction{
				table:           "filter",
				chain:           "INPUT",
				append:          true,
				inputInterface:  "tun0",
				protocol:        "tcp",
				source:          netip.MustParsePrefix("1.2.3.4/32"),
				destinationPort: 10000,
				target:          "ACCEPT",
			},
		},
		"nat_redirection": {
			s: "-t nat --delete PREROUTING -i tun0 -p tcp --dport 43716 -j REDIRECT --to-ports 5678",
			instruction: iptablesInstruction{
				table:           "nat",
				chain:           "PREROUTING",
				append:          false,
				inputInterface:  "tun0",
				protocol:        "tcp",
				destinationPort: 43716,
				target:          "REDIRECT",
				toPorts:         []uint16{5678},
			},
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			rule, err := parseIptablesInstruction(testCase.s)

			assert.Equal(t, testCase.instruction, rule)
			assert.ErrorIs(t, err, testCase.errWrapped)
			if testCase.errWrapped != nil {
				assert.EqualError(t, err, testCase.errMessage)
			}
		})
	}
}
