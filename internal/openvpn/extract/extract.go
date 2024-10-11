package extract

import (
	"errors"
	"fmt"
	"net/netip"
	"strconv"
	"strings"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
)

var errRemoteLineNotFound = errors.New("remote line not found")

func extractDataFromLines(lines []string) (
	connection models.Connection, err error,
) {
	for i, line := range lines {
		hashSymbolIndex := strings.Index(line, "#")
		if hashSymbolIndex >= 0 {
			line = line[:hashSymbolIndex]
		}

		ip, port, protocol, err := extractDataFromLine(line)
		if err != nil {
			return connection, fmt.Errorf("on line %d: %w", i+1, err)
		}

		connection.UpdateEmptyWith(ip, port, protocol)

		if connection.Protocol != "" && connection.IP.IsValid() {
			break
		}
	}

	if !connection.IP.IsValid() {
		return connection, errRemoteLineNotFound
	}

	if connection.Protocol == "" {
		connection.Protocol = constants.UDP
	}

	if connection.Port == 0 {
		connection.Port = 1194
		if strings.HasPrefix(connection.Protocol, "tcp") {
			connection.Port = 443
		}
	}

	return connection, nil
}

func extractDataFromLine(line string) (
	ip netip.Addr, port uint16, protocol string, err error,
) {
	switch {
	case strings.HasPrefix(line, "proto "):
		protocol, err = extractProto(line)
		if err != nil {
			return ip, 0, "", fmt.Errorf("extracting protocol from proto line: %w", err)
		}
		return ip, 0, protocol, nil

	case strings.HasPrefix(line, "remote "):
		ip, port, protocol, err = extractRemote(line)
		if err != nil {
			return ip, 0, "", fmt.Errorf("extracting from remote line: %w", err)
		}
		return ip, port, protocol, nil

	case strings.HasPrefix(line, "port "):
		port, err = extractPort(line)
		if err != nil {
			return ip, 0, "", fmt.Errorf("extracting from port line: %w", err)
		}
		return ip, port, "", nil
	}

	return ip, 0, "", nil
}

var (
	errProtoLineFieldsCount = errors.New("proto line has not 2 fields as expected")
	errProtocolNotSupported = errors.New("network protocol not supported")
)

func extractProto(line string) (protocol string, err error) {
	fields := strings.Fields(line)
	if len(fields) != 2 { //nolint:mnd
		return "", fmt.Errorf("%w: %s", errProtoLineFieldsCount, line)
	}

	switch fields[1] {
	case "tcp", "tcp4", "tcp6", "tcp-client", "udp", "udp4", "udp6":
	default:
		return "", fmt.Errorf("%w: %s", errProtocolNotSupported, fields[1])
	}

	return fields[1], nil
}

var (
	errRemoteLineFieldsCount = errors.New("remote line has not 2 fields as expected")
	errHostNotIP             = errors.New("host is not an IP address")
	errPortNotValid          = errors.New("port is not valid")
)

func extractRemote(line string) (ip netip.Addr, port uint16,
	protocol string, err error,
) {
	fields := strings.Fields(line)
	n := len(fields)

	if n < 2 || n > 4 {
		return netip.Addr{}, 0, "", fmt.Errorf("%w: %s", errRemoteLineFieldsCount, line)
	}

	host := fields[1]
	ip, err = netip.ParseAddr(host)
	if err != nil {
		return netip.Addr{}, 0, "", fmt.Errorf("%w: %s", errHostNotIP, host)
		// TODO resolve hostname once there is an option to allow it through
		// the firewall before the VPN is up.
	}

	if n > 2 { //nolint:mnd
		portInt, err := strconv.Atoi(fields[2])
		if err != nil {
			return netip.Addr{}, 0, "", fmt.Errorf("%w: %s", errPortNotValid, line)
		} else if portInt < 1 || portInt > 65535 {
			return netip.Addr{}, 0, "", fmt.Errorf("%w: %d must be between 1 and 65535", errPortNotValid, portInt)
		}
		port = uint16(portInt)
	}

	if n > 3 { //nolint:mnd
		switch fields[3] {
		case "tcp", "udp":
			protocol = fields[3]
		default:
			return netip.Addr{}, 0, "", fmt.Errorf("%w: %s", errProtocolNotSupported, fields[3])
		}
	}

	return ip, port, protocol, nil
}

var errPostLineFieldsCount = errors.New("post line has not 2 fields as expected")

func extractPort(line string) (port uint16, err error) {
	fields := strings.Fields(line)
	const expectedFieldsCount = 2
	if len(fields) != expectedFieldsCount {
		return 0, fmt.Errorf("%w: %s", errPostLineFieldsCount, line)
	}

	portInt, err := strconv.Atoi(fields[1])
	if err != nil {
		return 0, fmt.Errorf("%w: %s", errPortNotValid, line)
	} else if portInt < 1 || portInt > 65535 {
		return 0, fmt.Errorf("%w: %d must be between 1 and 65535", errPortNotValid, portInt)
	}
	port = uint16(portInt)

	return port, nil
}
