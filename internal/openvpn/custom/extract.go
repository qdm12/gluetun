package custom

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
	connection models.Connection, intf string, err error) {
	for i, line := range lines {
		ip, port, protocol, intfFound, err := extractDataFromLine(line)
		if err != nil {
			return connection, "", fmt.Errorf("on line %d: %w", i+1, err)
		}

		intf = intfFound
		connection.UpdateEmptyWith(ip, port, protocol)

		if connection.Protocol != "" && connection.IP != nil && intf != "" {
			break
		}
	}

	if connection.IP == nil {
		return connection, "", errRemoteLineNotFound
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

	return connection, intf, nil
}

var (
	errExtractProto  = errors.New("failed extracting protocol from proto line")
	errExtractRemote = errors.New("failed extracting from remote line")
	errExtractDev    = errors.New("failed extracting network interface from dev line")
)

func extractDataFromLine(line string) (
	ip net.IP, port uint16, protocol, intf string, err error) {
	switch {
	case strings.HasPrefix(line, "proto "):
		protocol, err = extractProto(line)
		if err != nil {
			return nil, 0, "", "", fmt.Errorf("%w: %s", errExtractProto, err)
		}
		return nil, 0, protocol, "", nil

	case strings.HasPrefix(line, "remote "):
		ip, port, protocol, err = extractRemote(line)
		if err != nil {
			return nil, 0, "", "", fmt.Errorf("%w: %s", errExtractRemote, err)
		}
		return ip, port, protocol, "", nil

	case strings.HasPrefix(line, "dev "):
		intf, err = extractInterfaceFromLine(line)
		if err != nil {
			return nil, 0, "", "", fmt.Errorf("%w: %s", errExtractDev, err)
		}
		return nil, 0, "", intf, nil
	}

	return nil, 0, "", "", nil
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
	case "tcp", "udp":
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
			return nil, 0, "", fmt.Errorf("%w: not between 1 and 65535: %d", errPortNotValid, portInt)
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

var (
	errDevLineFieldsCount = errors.New("dev line has not 2 fields as expected")
)

func extractInterfaceFromLine(line string) (intf string, err error) {
	fields := strings.Fields(line)
	if len(fields) != 2 { //nolint:gomnd
		return "", fmt.Errorf("%w: %s", errDevLineFieldsCount, line)
	}

	return fields[1], nil
}
