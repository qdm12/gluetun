package firewall

import (
	"errors"
	"fmt"
	"net/netip"
	"slices"
	"strconv"
	"strings"
)

type chain struct {
	name    string
	policy  string
	packets uint64
	bytes   uint64
	rules   []chainRule
}

type chainRule struct {
	lineNumber      uint16 // starts from 1 and cannot be zero.
	packets         uint64
	bytes           uint64
	target          string       // "ACCEPT", "DROP", "REJECT" or "REDIRECT"
	protocol        string       // "icmp", "tcp", "udp" or "" for all protocols.
	inputInterface  string       // input interface, for example "tun0" or "*""
	outputInterface string       // output interface, for example "eth0" or "*""
	source          netip.Prefix // source IP CIDR, for example 0.0.0.0/0. Must be valid.
	destination     netip.Prefix // destination IP CIDR, for example 0.0.0.0/0. Must be valid.
	destinationPort uint16       // Not specified if set to zero.
	redirPorts      []uint16     // Not specified if empty.
	ctstate         []string     // for example ["RELATED","ESTABLISHED"]. Can be empty.
}

var ErrChainListMalformed = errors.New("iptables chain list output is malformed")

func parseChain(iptablesOutput string) (c chain, err error) {
	// Text example:
	// Chain INPUT (policy ACCEPT 140K packets, 226M bytes)
	// pkts bytes target     prot opt in     out     source               destination
	// 	  0     0 ACCEPT     17   --  tun0   *       0.0.0.0/0            0.0.0.0/0            udp dpt:55405
	// 	  0     0 ACCEPT     6    --  tun0   *       0.0.0.0/0            0.0.0.0/0            tcp dpt:55405
	// 	  0     0 DROP       0    --  tun0   *       0.0.0.0/0            0.0.0.0/0
	iptablesOutput = strings.TrimSpace(iptablesOutput)
	linesWithComments := strings.Split(iptablesOutput, "\n")

	// Filter out lines starting with a '#' character
	lines := make([]string, 0, len(linesWithComments))
	for _, line := range linesWithComments {
		if strings.HasPrefix(line, "#") {
			continue
		}
		lines = append(lines, line)
	}

	const minLines = 2 // chain general information line + legend line
	if len(lines) < minLines {
		return chain{}, fmt.Errorf("%w: not enough lines to process in: %s",
			ErrChainListMalformed, iptablesOutput)
	}

	c, err = parseChainGeneralDataLine(lines[0])
	if err != nil {
		return chain{}, fmt.Errorf("parsing chain general data line: %w", err)
	}

	// Sanity check for the legend line
	expectedLegendFields := []string{"num", "pkts", "bytes", "target", "prot", "opt", "in", "out", "source", "destination"}
	legendLine := strings.TrimSpace(lines[1])
	legendFields := strings.Fields(legendLine)
	if !slices.Equal(expectedLegendFields, legendFields) {
		return chain{}, fmt.Errorf("%w: legend %q is not the expected %q",
			ErrChainListMalformed, legendLine, strings.Join(expectedLegendFields, " "))
	}

	lines = lines[2:] // remove chain general information line and legend line
	if len(lines) == 0 {
		return c, nil
	}

	c.rules = make([]chainRule, len(lines))
	for i, line := range lines {
		c.rules[i], err = parseChainRuleLine(line)
		if err != nil {
			return chain{}, fmt.Errorf("parsing chain rule %q: %w", line, err)
		}
	}

	return c, nil
}

// parseChainGeneralDataLine parses the first line of iptables chain list output.
// For example, it can parse the following line:
// Chain INPUT (policy ACCEPT 140K packets, 226M bytes)
// It returns a chain struct with the parsed data.
func parseChainGeneralDataLine(line string) (base chain, err error) {
	line = strings.TrimSpace(line)
	runesToRemove := []rune{'(', ')', ','}
	for _, r := range runesToRemove {
		line = strings.ReplaceAll(line, string(r), "")
	}

	fields := strings.Fields(line)
	const expectedNumberOfFields = 8
	if len(fields) != expectedNumberOfFields {
		return chain{}, fmt.Errorf("%w: expected %d fields in %q",
			ErrChainListMalformed, expectedNumberOfFields, line)
	}

	// Sanity checks
	indexToExpectedValue := map[int]string{
		0: "Chain",
		2: "policy",
		5: "packets",
		7: "bytes",
	}
	for index, expectedValue := range indexToExpectedValue {
		if fields[index] == expectedValue {
			continue
		}
		return chain{}, fmt.Errorf("%w: expected %q for field %d in %q",
			ErrChainListMalformed, expectedValue, index, line)
	}

	base.name = fields[1] // chain name could be custom
	base.policy = fields[3]
	err = checkTarget(base.policy)
	if err != nil {
		return chain{}, fmt.Errorf("policy target in %q: %w", line, err)
	}

	packets, err := parseMetricSize(fields[4])
	if err != nil {
		return chain{}, fmt.Errorf("parsing packets: %w", err)
	}
	base.packets = packets

	bytes, err := parseMetricSize(fields[6])
	if err != nil {
		return chain{}, fmt.Errorf("parsing bytes: %w", err)
	}
	base.bytes = bytes

	return base, nil
}

var ErrChainRuleMalformed = errors.New("chain rule is malformed")

func parseChainRuleLine(line string) (rule chainRule, err error) {
	line = strings.TrimSpace(line)
	if line == "" {
		return chainRule{}, fmt.Errorf("%w: empty line", ErrChainRuleMalformed)
	}

	fields := strings.Fields(line)

	const minFields = 10
	if len(fields) < minFields {
		return chainRule{}, fmt.Errorf("%w: not enough fields", ErrChainRuleMalformed)
	}

	for fieldIndex, field := range fields[:minFields] {
		err = parseChainRuleField(fieldIndex, field, &rule)
		if err != nil {
			return chainRule{}, fmt.Errorf("parsing chain rule field: %w", err)
		}
	}

	if len(fields) > minFields {
		err = parseChainRuleOptionalFields(fields[minFields:], &rule)
		if err != nil {
			return chainRule{}, fmt.Errorf("parsing optional fields: %w", err)
		}
	}

	return rule, nil
}

func parseChainRuleField(fieldIndex int, field string, rule *chainRule) (err error) {
	if field == "" {
		return fmt.Errorf("%w: empty field at index %d", ErrChainRuleMalformed, fieldIndex)
	}

	const (
		numIndex = iota
		packetsIndex
		bytesIndex
		targetIndex
		protocolIndex
		optIndex
		inputInterfaceIndex
		outputInterfaceIndex
		sourceIndex
		destinationIndex
	)

	switch fieldIndex {
	case numIndex:
		rule.lineNumber, err = parseLineNumber(field)
		if err != nil {
			return fmt.Errorf("parsing line number: %w", err)
		}
	case packetsIndex:
		rule.packets, err = parseMetricSize(field)
		if err != nil {
			return fmt.Errorf("parsing packets: %w", err)
		}
	case bytesIndex:
		rule.bytes, err = parseMetricSize(field)
		if err != nil {
			return fmt.Errorf("parsing bytes: %w", err)
		}
	case targetIndex:
		err = checkTarget(field)
		if err != nil {
			return fmt.Errorf("checking target: %w", err)
		}
		rule.target = field
	case protocolIndex:
		rule.protocol, err = parseProtocol(field)
		if err != nil {
			return fmt.Errorf("parsing protocol: %w", err)
		}
	case optIndex: // ignored
	case inputInterfaceIndex:
		rule.inputInterface = field
	case outputInterfaceIndex:
		rule.outputInterface = field
	case sourceIndex:
		rule.source, err = parseIPPrefix(field)
		if err != nil {
			return fmt.Errorf("parsing source IP CIDR: %w", err)
		}
	case destinationIndex:
		rule.destination, err = parseIPPrefix(field)
		if err != nil {
			return fmt.Errorf("parsing destination IP CIDR: %w", err)
		}
	}
	return nil
}

func parseChainRuleOptionalFields(optionalFields []string, rule *chainRule) (err error) {
	for i := 0; i < len(optionalFields); i++ {
		key := optionalFields[i]
		switch key {
		case "tcp", "udp":
			i++
			value := optionalFields[i]
			value = strings.TrimPrefix(value, "dpt:")
			const base, bitLength = 10, 16
			destinationPort, err := strconv.ParseUint(value, base, bitLength)
			if err != nil {
				return fmt.Errorf("parsing destination port %q: %w", value, err)
			}
			rule.destinationPort = uint16(destinationPort)
		case "redir":
			i++
			switch optionalFields[i] {
			case "ports":
				i++
				ports, err := parsePortsCSV(optionalFields[i])
				if err != nil {
					return fmt.Errorf("parsing redirection ports: %w", err)
				}
				rule.redirPorts = ports
			default:
				return fmt.Errorf("%w: unexpected optional field: %s",
					ErrChainRuleMalformed, optionalFields[i])
			}
		case "ctstate":
			i++
			rule.ctstate = strings.Split(optionalFields[i], ",")
		default:
			return fmt.Errorf("%w: unexpected optional field: %s", ErrChainRuleMalformed, key)
		}
	}
	return nil
}

func parsePortsCSV(s string) (ports []uint16, err error) {
	if s == "" {
		return nil, nil
	}

	fields := strings.Split(s, ",")
	ports = make([]uint16, len(fields))
	for i, field := range fields {
		const base, bitLength = 10, 16
		port, err := strconv.ParseUint(field, base, bitLength)
		if err != nil {
			return nil, fmt.Errorf("parsing port %q: %w", field, err)
		}
		ports[i] = uint16(port)
	}
	return ports, nil
}

var ErrLineNumberIsZero = errors.New("line number is zero")

func parseLineNumber(s string) (n uint16, err error) {
	const base, bitLength = 10, 16
	lineNumber, err := strconv.ParseUint(s, base, bitLength)
	if err != nil {
		return 0, err
	} else if lineNumber == 0 {
		return 0, fmt.Errorf("%w", ErrLineNumberIsZero)
	}
	return uint16(lineNumber), nil
}

var ErrTargetUnknown = errors.New("unknown target")

func checkTarget(target string) (err error) {
	switch target {
	case "ACCEPT", "DROP", "REJECT", "REDIRECT":
		return nil
	}
	return fmt.Errorf("%w: %s", ErrTargetUnknown, target)
}

var ErrProtocolUnknown = errors.New("unknown protocol")

func parseProtocol(s string) (protocol string, err error) {
	switch s {
	case "0":
	case "1":
		protocol = "icmp"
	case "6":
		protocol = "tcp"
	case "17":
		protocol = "udp"
	default:
		return "", fmt.Errorf("%w: %s", ErrProtocolUnknown, s)
	}
	return protocol, nil
}

var ErrMetricSizeMalformed = errors.New("metric size is malformed")

// parseMetricSize parses a metric size string like 140K or 226M and
// returns the raw integer matching it.
func parseMetricSize(size string) (n uint64, err error) {
	if size == "" {
		return n, fmt.Errorf("%w: empty string", ErrMetricSizeMalformed)
	}

	//nolint:mnd
	multiplerLetterToValue := map[byte]uint64{
		'K': 1000,
		'M': 1000000,
		'G': 1000000000,
		'T': 1000000000000,
	}

	lastCharacter := size[len(size)-1]
	multiplier, ok := multiplerLetterToValue[lastCharacter]
	if ok { // multiplier present
		size = size[:len(size)-1]
	} else {
		multiplier = 1
	}

	const base, bitLength = 10, 64
	n, err = strconv.ParseUint(size, base, bitLength)
	if err != nil {
		return n, fmt.Errorf("%w: %w", ErrMetricSizeMalformed, err)
	}
	n *= multiplier
	return n, nil
}
