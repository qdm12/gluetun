package custom

import (
	"errors"
	"net"
	"testing"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_extractDataFromLines(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		lines      []string
		connection models.Connection
		intf       string
		err        error
	}{
		"success": {
			lines: []string{"bla bla", "proto tcp", "remote 1.2.3.4 1194 tcp", "dev tun6"},
			connection: models.Connection{
				IP:       net.IPv4(1, 2, 3, 4),
				Port:     1194,
				Protocol: constants.TCP,
			},
			intf: "tun6",
		},
		"extraction error": {
			lines: []string{"bla bla", "proto bad", "remote 1.2.3.4 1194 tcp"},
			err:   errors.New("on line 2: failed extracting protocol from proto line: network protocol not supported: bad"),
		},
		"only use first values found": {
			lines: []string{"proto udp", "proto tcp", "remote 1.2.3.4 443 tcp", "remote 5.2.3.4 1194 udp"},
			connection: models.Connection{
				IP:       net.IPv4(1, 2, 3, 4),
				Port:     443,
				Protocol: constants.UDP,
			},
		},
		"no IP found": {
			lines: []string{"proto tcp"},
			connection: models.Connection{
				Protocol: constants.TCP,
			},
			err: errRemoteLineNotFound,
		},
		"default TCP port": {
			lines: []string{"remote 1.2.3.4", "proto tcp"},
			connection: models.Connection{
				IP:       net.IPv4(1, 2, 3, 4),
				Port:     443,
				Protocol: constants.TCP,
			},
		},
		"default UDP port": {
			lines: []string{"remote 1.2.3.4", "proto udp"},
			connection: models.Connection{
				IP:       net.IPv4(1, 2, 3, 4),
				Port:     1194,
				Protocol: constants.UDP,
			},
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			connection, intf, err := extractDataFromLines(testCase.lines)

			if testCase.err != nil {
				require.Error(t, err)
				assert.Equal(t, testCase.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, testCase.connection, connection)
			assert.Equal(t, testCase.intf, intf)
		})
	}
}

func Test_extractDataFromLine(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		line     string
		ip       net.IP
		port     uint16
		protocol string
		intf     string
		isErr    error
	}{
		"irrelevant line": {
			line: "bla bla",
		},
		"extract proto error": {
			line:  "proto bad",
			isErr: errExtractProto,
		},
		"extract proto success": {
			line:     "proto tcp",
			protocol: constants.TCP,
		},
		"extract intf error": {
			line:  "dev ",
			isErr: errExtractDev,
		},
		"extract intf success": {
			line: "dev tun3",
			intf: "tun3",
		},
		"extract remote error": {
			line:  "remote bad",
			isErr: errExtractRemote,
		},
		"extract remote success": {
			line:     "remote 1.2.3.4 1194 udp",
			ip:       net.IPv4(1, 2, 3, 4),
			port:     1194,
			protocol: constants.UDP,
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ip, port, protocol, intf, err := extractDataFromLine(testCase.line)

			if testCase.isErr != nil {
				assert.ErrorIs(t, err, testCase.isErr)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, testCase.ip, ip)
			assert.Equal(t, testCase.port, port)
			assert.Equal(t, testCase.protocol, protocol)
			assert.Equal(t, testCase.intf, intf)
		})
	}
}

func Test_extractProto(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		line     string
		protocol string
		err      error
	}{
		"fields error": {
			line: "proto one two",
			err:  errors.New("proto line has not 2 fields as expected: proto one two"),
		},
		"bad protocol": {
			line: "proto bad",
			err:  errors.New("network protocol not supported: bad"),
		},
		"udp": {
			line:     "proto udp",
			protocol: constants.UDP,
		},
		"tcp": {
			line:     "proto tcp",
			protocol: constants.TCP,
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			protocol, err := extractProto(testCase.line)

			if testCase.err != nil {
				require.Error(t, err)
				assert.Equal(t, testCase.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, testCase.protocol, protocol)
		})
	}
}

func Test_extractRemote(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		line     string
		ip       net.IP
		port     uint16
		protocol string
		err      error
	}{
		"not enough fields": {
			line: "remote",
			err:  errors.New("remote line has not 2 fields as expected: remote"),
		},
		"too many fields": {
			line: "remote one two three four",
			err:  errors.New("remote line has not 2 fields as expected: remote one two three four"),
		},
		"host is not an IP": {
			line: "remote somehost.com",
			err:  errors.New("host is not an an IP address: somehost.com"),
		},
		"only IP host": {
			line: "remote 1.2.3.4",
			ip:   net.IPv4(1, 2, 3, 4),
		},
		"port not an integer": {
			line: "remote 1.2.3.4 bad",
			err:  errors.New("port is not valid: remote 1.2.3.4 bad"),
		},
		"port is zero": {
			line: "remote 1.2.3.4 0",
			err:  errors.New("port is not valid: not between 1 and 65535: 0"),
		},
		"port is minus one": {
			line: "remote 1.2.3.4 -1",
			err:  errors.New("port is not valid: not between 1 and 65535: -1"),
		},
		"port is over 65535": {
			line: "remote 1.2.3.4 65536",
			err:  errors.New("port is not valid: not between 1 and 65535: 65536"),
		},
		"IP host and port": {
			line: "remote 1.2.3.4 8000",
			ip:   net.IPv4(1, 2, 3, 4),
			port: 8000,
		},
		"invalid protocol": {
			line: "remote 1.2.3.4 8000 bad",
			err:  errors.New("network protocol not supported: bad"),
		},
		"IP host and port and protocol": {
			line:     "remote 1.2.3.4 8000 udp",
			ip:       net.IPv4(1, 2, 3, 4),
			port:     8000,
			protocol: constants.UDP,
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ip, port, protocol, err := extractRemote(testCase.line)

			if testCase.err != nil {
				require.Error(t, err)
				assert.Equal(t, testCase.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, testCase.ip, ip)
			assert.Equal(t, testCase.port, port)
			assert.Equal(t, testCase.protocol, protocol)
		})
	}
}

func Test_extractInterface(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		line string
		intf string
		err  error
	}{
		"found": {
			line: "dev tun3",
			intf: "tun3",
		},
		"not enough fields": {
			line: "dev ",
			err:  errors.New("dev line has not 2 fields as expected: dev "),
		},
		"too many fields": {
			line: "dev one two",
			err:  errors.New("dev line has not 2 fields as expected: dev one two"),
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			intf, err := extractInterfaceFromLine(testCase.line)

			if testCase.err != nil {
				require.Error(t, err)
				assert.Equal(t, testCase.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, testCase.intf, intf)
		})
	}
}
