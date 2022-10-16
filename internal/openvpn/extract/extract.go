package extract

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
)

var (
	errRemoteLineNotFound = errors.New("remote line not found")
)

func extractDataFromLines(lines []string) (
	connection models.Connection, err error) {
	for i, line := range lines {
		ip, port, protocol, err := extractDataFromLine(line)
		if err != nil {
			return connection, fmt.Errorf("on line %d: %w", i+1, err)
		}

		connection.UpdateEmptyWith(ip, port, protocol)

		if connection.Protocol != "" && connection.IP != nil {
			break
		}
	}

	if connection.IP == nil {
		return connection, errRemoteLineNotFound
	}

	if connection.Protocol == "" {
		connection.Protocol = constants.UDP
	}

	if connection.Port == 0 {
		connection.Port = 1194
		if connection.Protocol == constants.TCP {
			connection.Port = 443
		}
	}

	return connection, nil
}

func extractDataFromLine(line string) (
	ip net.IP, port uint16, protocol string, err error) {
	switch {
	case strings.HasPrefix(line, "proto "):
		protocol, err = extractProto(line)
		if err != nil {
			return nil, 0, "", fmt.Errorf("failed extracting protocol from proto line: %w", err)
		}
		return nil, 0, protocol, nil

	case strings.HasPrefix(line, "remote "):
		ip, port, protocol, err = extractRemote(line)
		if err != nil {
			return nil, 0, "", fmt.Errorf("failed extracting from remote line: %w", err)
		}
		return ip, port, protocol, nil
	}

	return nil, 0, "", nil
}

var (
	errProtoLineFieldsCount = errors.New("proto line has not 2 fields as expected")
	errProtocolNotSupported = errors.New("network protocol not supported")
)

func extractProto(line string) (protocol string, err error) {
	fields := strings.Fields(line)
	if len(fields) != 2 { //nolint:gomnd
		return "", fmt.Errorf("%w: %s", errProtoLineFieldsCount, line)
	}

	switch fields[1] {
	case "tcp", "tcp4", "tcp6", "udp", "udp4", "udp6":
	default:
		return "", fmt.Errorf("%w: %s", errProtocolNotSupported, fields[1])
	}

	return fields[1], nil
}

var (
	errRemoteLineFieldsCount = errors.New("remote line has not 2 fields as expected")
	errHostNotIP             = errors.New("host is not an an IP address")
	errPortNotValid          = errors.New("port is not valid")
)

func extractRemote(line string) (ip net.IP, port uint16,
	protocol string, err error) {
	fields := strings.Fields(line)
	n := len(fields)

	if n < 2 || n > 4 {
		return nil, 0, "", fmt.Errorf("%w: %s", errRemoteLineFieldsCount, line)
	}

	host := fields[1]
	ip = net.ParseIP(host)
	if ip == nil {
		return nil, 0, "", fmt.Errorf("%w: %s", errHostNotIP, host)
		// TODO resolve hostname once there is an option to allow it through
		// the firewall before the VPN is up.
	}

	if n > 2 { //nolint:gomnd
		portInt, err := strconv.Atoi(fields[2])
		if err != nil {
			return nil, 0, "", fmt.Errorf("%w: %s", errPortNotValid, line)
		} else if portInt < 1 || portInt > 65535 {
			return nil, 0, "", fmt.Errorf("%w: %d must be between 1 and 65535", errPortNotValid, portInt)
		}
		port = uint16(portInt)
	}

	if n > 3 { //nolint:gomnd
		switch fields[3] {
		case "tcp", "udp":
			protocol = fields[3]
		default:
			return nil, 0, "", fmt.Errorf("%w: %s", errProtocolNotSupported, fields[3])
		}
	}

	return ip, port, protocol, nil
}
