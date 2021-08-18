package openvpn

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

var (
	errReadCustomConfig  = errors.New("cannot read custom configuration file")
	errExtractConnection = errors.New("cannot extract connection from custom configuration file")
)

func processCustomConfig(settings configuration.OpenVPN) (
	lines []string, connection models.OpenVPNConnection, err error) {
	lines, err = readCustomConfigLines(settings.Config)
	if err != nil {
		return nil, connection, fmt.Errorf("%w: %s", errReadCustomConfig, err)
	}

	connection, err = extractConnectionFromLines(lines)
	if err != nil {
		return nil, connection, fmt.Errorf("%w: %s", errExtractConnection, err)
	}

	lines = modifyCustomConfig(lines, settings, connection)

	return lines, connection, nil
}

func readCustomConfigLines(filepath string) (
	lines []string, err error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	b, err := io.ReadAll(file)
	if err != nil {
		_ = file.Close()
		return nil, err
	}

	if err := file.Close(); err != nil {
		return nil, err
	}

	return strings.Split(string(b), "\n"), nil
}

func modifyCustomConfig(lines []string, settings configuration.OpenVPN,
	connection models.OpenVPNConnection) (modified []string) {
	// Remove some lines
	for _, line := range lines {
		switch {
		case strings.HasPrefix(line, "up "),
			strings.HasPrefix(line, "down "),
			strings.HasPrefix(line, "verb "),
			strings.HasPrefix(line, "auth-user-pass "),
			strings.HasPrefix(line, "user "),
			strings.HasPrefix(line, "proto "),
			strings.HasPrefix(line, "remote "),
			settings.Cipher != "" && strings.HasPrefix(line, "cipher "),
			settings.Cipher != "" && strings.HasPrefix(line, "data-ciphers "),
			settings.Auth != "" && strings.HasPrefix(line, "auth "),
			settings.MSSFix > 0 && strings.HasPrefix(line, "mssfix "),
			!settings.IPv6 && strings.HasPrefix(line, "tun-ipv6"):
		default:
			modified = append(modified, line)
		}
	}

	// Add values
	modified = append(modified, connection.ProtoLine())
	modified = append(modified, connection.RemoteLine())
	modified = append(modified, "mute-replay-warnings")
	modified = append(modified, "auth-nocache")
	modified = append(modified, "pull-filter ignore \"auth-token\"") // prevent auth failed loop
	modified = append(modified, "auth-retry nointeract")
	modified = append(modified, "suppress-timestamps")
	if settings.User != "" {
		modified = append(modified, "auth-user-pass "+constants.OpenVPNAuthConf)
	}
	modified = append(modified, "verb "+strconv.Itoa(settings.Verbosity))
	if settings.Cipher != "" {
		modified = append(modified, utils.CipherLines(settings.Cipher, settings.Version)...)
	}
	if settings.Auth != "" {
		modified = append(modified, "auth "+settings.Auth)
	}
	if settings.MSSFix > 0 {
		modified = append(modified, "mssfix "+strconv.Itoa(int(settings.MSSFix)))
	}
	if !settings.IPv6 {
		modified = append(modified, `pull-filter ignore "route-ipv6"`)
		modified = append(modified, `pull-filter ignore "ifconfig-ipv6"`)
	}
	if !settings.Root {
		modified = append(modified, "user "+settings.ProcUser)
	}

	return modified
}

var (
	errRemoteLineNotFound = errors.New("remote line not found")
)

// extractConnectionFromLines always takes the first remote line only.
func extractConnectionFromLines(lines []string) (
	connection models.OpenVPNConnection, err error) {
	for i, line := range lines {
		newConnectionData, err := extractConnectionFromLine(line)
		if err != nil {
			return connection, fmt.Errorf("on line %d: %w", i+1, err)
		}
		connection.UpdateEmptyWith(newConnectionData)

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

var (
	errExtractProto  = errors.New("failed extracting protocol from proto line")
	errExtractRemote = errors.New("failed extracting protocol from remote line")
)

func extractConnectionFromLine(line string) (
	connection models.OpenVPNConnection, err error) {
	switch {
	case strings.HasPrefix(line, "proto "):
		connection.Protocol, err = extractProto(line)
		if err != nil {
			return connection, fmt.Errorf("%w: %s", errExtractProto, err)
		}

	// only take the first remote line
	case strings.HasPrefix(line, "remote ") && connection.IP == nil:
		connection.IP, connection.Port, connection.Protocol, err = extractRemote(line)
		if err != nil {
			return connection, fmt.Errorf("%w: %s", errExtractRemote, err)
		}
	}

	return connection, nil
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
