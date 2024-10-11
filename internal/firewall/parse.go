package firewall

import (
	"errors"
	"fmt"
	"net/netip"
	"slices"
	"strconv"
	"strings"
)

type iptablesInstruction struct {
	table           string // defaults to "filter", and can be "nat" for example.
	append          bool
	chain           string       // for example INPUT, PREROUTING. Cannot be empty.
	target          string       // for example ACCEPT. Can be empty.
	protocol        string       // "tcp" or "udp" or "" for all protocols.
	inputInterface  string       // for example "tun0" or "" for any interface.
	outputInterface string       // for example "tun0" or "" for any interface.
	source          netip.Prefix // if not valid, then it is unspecified.
	destination     netip.Prefix // if not valid, then it is unspecified.
	destinationPort uint16       // if zero, there is no destination port
	toPorts         []uint16     // if empty, there is no redirection
	ctstate         []string     // if empty, there is no ctstate
}

func (i *iptablesInstruction) setDefaults() {
	if i.table == "" {
		i.table = "filter"
	}
}

// equalToRule ignores the append boolean flag of the instruction to compare against the rule.
func (i *iptablesInstruction) equalToRule(table, chain string, rule chainRule) (equal bool) {
	switch {
	case i.table != table:
		return false
	case i.chain != chain:
		return false
	case i.target != rule.target:
		return false
	case i.protocol != rule.protocol:
		return false
	case i.destinationPort != rule.destinationPort:
		return false
	case !slices.Equal(i.toPorts, rule.redirPorts):
		return false
	case !slices.Equal(i.ctstate, rule.ctstate):
		return false
	case !networkInterfacesEqual(i.inputInterface, rule.inputInterface):
		return false
	case !networkInterfacesEqual(i.outputInterface, rule.outputInterface):
		return false
	case !ipPrefixesEqual(i.source, rule.source):
		return false
	case !ipPrefixesEqual(i.destination, rule.destination):
		return false
	default:
		return true
	}
}

// instruction can be "" which equivalent to the "*" chain rule interface.
func networkInterfacesEqual(instruction, chainRule string) bool {
	return instruction == chainRule || (instruction == "" && chainRule == "*")
}

func ipPrefixesEqual(instruction, chainRule netip.Prefix) bool {
	return instruction == chainRule ||
		(!instruction.IsValid() && chainRule.Bits() == 0 && chainRule.Addr().IsUnspecified())
}

var ErrIptablesCommandMalformed = errors.New("iptables command is malformed")

func parseIptablesInstruction(s string) (instruction iptablesInstruction, err error) {
	if s == "" {
		return iptablesInstruction{}, fmt.Errorf("%w: empty instruction", ErrIptablesCommandMalformed)
	}
	fields := strings.Fields(s)
	if len(fields)%2 != 0 {
		return iptablesInstruction{}, fmt.Errorf("%w: fields count %d is not even: %q",
			ErrIptablesCommandMalformed, len(fields), s)
	}

	for i := 0; i < len(fields); i += 2 {
		key := fields[i]
		value := fields[i+1]
		err = parseInstructionFlag(key, value, &instruction)
		if err != nil {
			return iptablesInstruction{}, fmt.Errorf("parsing %q: %w", s, err)
		}
	}

	instruction.setDefaults()
	return instruction, nil
}

func parseInstructionFlag(key, value string, instruction *iptablesInstruction) (err error) {
	switch key {
	case "-t", "--table":
		instruction.table = value
	case "-D", "--delete":
		instruction.append = false
		instruction.chain = value
	case "-A", "--append":
		instruction.append = true
		instruction.chain = value
	case "-j", "--jump":
		instruction.target = value
	case "-p", "--protocol":
		instruction.protocol = value
	case "-m", "--match": // ignore match
	case "-i", "--in-interface":
		instruction.inputInterface = value
	case "-o", "--out-interface":
		instruction.outputInterface = value
	case "-s", "--source":
		instruction.source, err = parseIPPrefix(value)
		if err != nil {
			return fmt.Errorf("parsing source IP CIDR: %w", err)
		}
	case "-d", "--destination":
		instruction.destination, err = parseIPPrefix(value)
		if err != nil {
			return fmt.Errorf("parsing destination IP CIDR: %w", err)
		}
	case "--dport":
		const base, bitLength = 10, 16
		destinationPort, err := strconv.ParseUint(value, base, bitLength)
		if err != nil {
			return fmt.Errorf("parsing destination port: %w", err)
		}
		instruction.destinationPort = uint16(destinationPort)
	case "--ctstate":
		instruction.ctstate = strings.Split(value, ",")
	case "--to-ports":
		portStrings := strings.Split(value, ",")
		instruction.toPorts = make([]uint16, len(portStrings))
		for i, portString := range portStrings {
			const base, bitLength = 10, 16
			port, err := strconv.ParseUint(portString, base, bitLength)
			if err != nil {
				return fmt.Errorf("parsing port redirection: %w", err)
			}
			instruction.toPorts[i] = uint16(port)
		}
	default:
		return fmt.Errorf("%w: unknown key %q", ErrIptablesCommandMalformed, key)
	}
	return nil
}

func parseIPPrefix(value string) (prefix netip.Prefix, err error) {
	slashIndex := strings.Index(value, "/")
	if slashIndex >= 0 {
		return netip.ParsePrefix(value)
	}

	ip, err := netip.ParseAddr(value)
	if err != nil {
		return netip.Prefix{}, fmt.Errorf("parsing IP address: %w", err)
	}
	return netip.PrefixFrom(ip, ip.BitLen()), nil
}
