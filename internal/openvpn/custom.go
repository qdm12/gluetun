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
	modified = append(modified, `pull-filter ignore "ping-restart"`)
	modified = append(modified, "auth-retry nointeract")
	modified = append(modified, "suppress-timestamps")
	modified = append(modified, "auth-user-pass "+constants.OpenVPNAuthConf)
	modified = append(modified, "verb "+strconv.Itoa(settings.Verbosity))
	if len(settings.Cipher) > 0 {
		modified = append(modified, "cipher "+settings.Cipher)
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
func extractConnectionFromLines(lines []string) (
	connection models.OpenVPNConnection, err error) {
	var foundProto, foundRemote bool

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
			foundProto = true

		case strings.HasPrefix(line, "remote "):
			fields := strings.Fields(line)
			n := len(fields)
			switch n {
			case 3: //nolint:gomnd
			case 4: //nolint:gomnd
				connection.Protocol = fields[3]
				foundProto = true
			default:
				return connection, fmt.Errorf(
					"%w: remote line has %d fields instead of 3 or 4: %s",
					errExtractConnection, n, line)
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

			port, err := strconv.Atoi(fields[2])
			if err != nil {
				return connection, fmt.Errorf(
					"%w: remote line has an invalid port: %s",
					errExtractConnection, line)
			}
			connection.Port = uint16(port)

			foundRemote = true
		}

		if foundProto && foundRemote {
			break
		}
	}

	if !foundProto {
		return connection, fmt.Errorf("%w: proto line not found", errExtractConnection)
	} else if !foundRemote {
		return connection, fmt.Errorf("%w: remote line not found", errExtractConnection)
	}

	return connection, nil
}

func setConnectionToLines(lines []string, connection models.OpenVPNConnection) (modified []string) {
	for i, line := range lines {
		switch {
		case strings.HasPrefix(line, "proto "):
			lines[i] = "proto " + connection.Protocol

		case strings.HasPrefix(line, "remote "):
			lines[i] = "remote " + connection.IP.String() + " " + strconv.Itoa(int(connection.Port))
		}
	}

	return lines
}
