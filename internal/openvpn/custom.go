package openvpn

import (
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"

	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
	"github.com/qdm12/golibs/os"
)

var errProcessCustomConfig = errors.New("cannot process custom config")

func (l *looper) processCustomConfig(settings configuration.OpenVPN) (
	lines []string, connection models.OpenVPNConnection, err error) {
	lines, err = readCustomConfigLines(settings.Config, l.openFile)
	if err != nil {
		return nil, connection, fmt.Errorf("%w: %s", errProcessCustomConfig, err)
	}

	lines = modifyCustomConfig(lines, l.username, settings)

	connection, err = extractConnectionFromLines(lines)
	if err != nil {
		return nil, connection, fmt.Errorf("%w: %s", errProcessCustomConfig, err)
	}

	lines = setConnectionToLines(lines, connection)
	return lines, connection, nil
}

func readCustomConfigLines(filepath string, openFile os.OpenFileFunc) (
	lines []string, err error) {
	file, err := openFile(filepath, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	b, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	if err := file.Close(); err != nil {
		return nil, err
	}

	return strings.Split(string(b), "\n"), nil
}

func modifyCustomConfig(lines []string, username string,
	settings configuration.OpenVPN) (modified []string) {
	// Remove some lines
	for _, line := range lines {
		switch {
		case strings.HasPrefix(line, "up "),
			strings.HasPrefix(line, "down "),
			strings.HasPrefix(line, "verb "),
			strings.HasPrefix(line, "auth-user-pass "),
			len(settings.Cipher) > 0 && strings.HasPrefix(line, "cipher "),
			len(settings.Cipher) > 0 && strings.HasPrefix(line, "data-ciphers"),
			len(settings.Auth) > 0 && strings.HasPrefix(line, "auth "),
			settings.MSSFix > 0 && strings.HasPrefix(line, "mssfix "),
			!settings.Provider.ExtraConfigOptions.OpenVPNIPv6 && strings.HasPrefix(line, "tun-ipv6"):
		default:
			modified = append(modified, line)
		}
	}

	// Add values
	modified = append(modified, "mute-replay-warnings")
	modified = append(modified, "auth-nocache")
	modified = append(modified, "pull-filter ignore \"auth-token\"") // prevent auth failed loop
	modified = append(modified, "auth-retry nointeract")
	modified = append(modified, "suppress-timestamps")
	modified = append(modified, "auth-user-pass "+constants.OpenVPNAuthConf)
	modified = append(modified, "verb "+strconv.Itoa(settings.Verbosity))
	if len(settings.Cipher) > 0 {
		modified = append(modified, utils.CipherLines(settings.Cipher, settings.Version)...)
	}
	if len(settings.Auth) > 0 {
		modified = append(modified, "auth "+settings.Auth)
	}
	if settings.MSSFix > 0 {
		modified = append(modified, "mssfix "+strconv.Itoa(int(settings.MSSFix)))
	}
	if !settings.Provider.ExtraConfigOptions.OpenVPNIPv6 {
		modified = append(modified, `pull-filter ignore "route-ipv6"`)
		modified = append(modified, `pull-filter ignore "ifconfig-ipv6"`)
	}
	if !settings.Root {
		modified = append(modified, "user "+username)
	}

	return modified
}

var errExtractConnection = errors.New("cannot extract connection")

// extractConnectionFromLines always takes the first remote line only.
func extractConnectionFromLines(lines []string) ( //nolint:gocognit
	connection models.OpenVPNConnection, err error) {
	for _, line := range lines {
		switch {
		case strings.HasPrefix(line, "proto "):
			fields := strings.Fields(line)
			if n := len(fields); n != 2 { //nolint:gomnd
				return connection, fmt.Errorf(
					"%w: proto line has %d fields instead of 2: %s",
					errExtractConnection, n, line)
			}
			connection.Protocol = fields[1]

		// only take the first remote line
		case strings.HasPrefix(line, "remote ") && connection.IP == nil:
			fields := strings.Fields(line)
			n := len(fields)
			//nolint:gomnd
			if n < 2 {
				return connection, fmt.Errorf(
					"%w: remote line has not enough fields: %s",
					errExtractConnection, line)
			}

			host := fields[1]
			if ip := net.ParseIP(host); ip != nil {
				connection.IP = ip
			} else {
				return connection, fmt.Errorf(
					"%w: for now, the remote line must contain an IP adddress: %s",
					errExtractConnection, line)
				// TODO resolve hostname once there is an option to allow it through
				// the firewall before the VPN is up.
			}

			if n > 2 { //nolint:gomnd
				port, err := strconv.Atoi(fields[2])
				if err != nil {
					return connection, fmt.Errorf(
						"%w: remote line has an invalid port: %s",
						errExtractConnection, line)
				}
				connection.Port = uint16(port)
			}

			if n > 3 { //nolint:gomnd
				connection.Protocol = strings.ToLower(fields[3])
			}

			if n > 4 { //nolint:gomnd
				return connection, fmt.Errorf(
					"%w: remote line has too many fields: %s",
					errExtractConnection, line)
			}
		}

		if connection.Protocol != "" && connection.IP != nil {
			break
		}
	}

	if connection.IP == nil {
		return connection, fmt.Errorf("%w: remote line not found", errExtractConnection)
	}

	switch connection.Protocol {
	case "":
		connection.Protocol = "udp"
	case "tcp", "udp":
	default:
		return connection, fmt.Errorf("%w: network protocol not supported: %s", errExtractConnection, connection.Protocol)
	}

	if connection.Port == 0 {
		if connection.Protocol == "tcp" {
			const defaultPort uint16 = 443
			connection.Port = defaultPort
		} else {
			const defaultPort uint16 = 1194
			connection.Port = defaultPort
		}
	}

	return connection, nil
}

func setConnectionToLines(lines []string, connection models.OpenVPNConnection) (modified []string) {
	for i, line := range lines {
		switch {
		case strings.HasPrefix(line, "proto "):
			lines[i] = connection.ProtoLine()

		case strings.HasPrefix(line, "remote "):
			lines[i] = connection.RemoteLine()
		}
	}

	return lines
}
