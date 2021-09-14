package custom

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/utils"
)

var ErrExtractData = errors.New("failed extracting information from custom configuration file")

func (p *Provider) BuildConf(connection models.Connection,
	settings configuration.OpenVPN) (lines []string, err error) {
	lines, _, err = p.extractor.Data(settings.ConfFile)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrExtractData, err)
	}

	lines = modifyConfig(lines, connection, settings)

	return lines, nil
}

func modifyConfig(lines []string, connection models.Connection,
	settings configuration.OpenVPN) (modified []string) {
	// Remove some lines
	for _, line := range lines {
		switch {
		case
			// Remove empty lines
			line == "",
			// Remove future to be duplicates
			line == "mute-replay-warnings",
			line == "auth-nocache",
			line == "pull-filter ignore \"auth-token\"",
			line == "auth-retry nointeract",
			line == "suppress-timestamps",
			line == "persist-tun",
			line == "persist-key",
			// Remove values always modified
			strings.HasPrefix(line, "verb "),
			strings.HasPrefix(line, "auth-user-pass "),
			strings.HasPrefix(line, "user "),
			strings.HasPrefix(line, "proto "),
			strings.HasPrefix(line, "remote "),
			strings.HasPrefix(line, "dev "),
			// Remove values eventually modified
			settings.Cipher != "" && hasPrefixOneOf(line,
				"cipher ", "data-ciphers ", "data-ciphers-fallback "),
			settings.Auth != "" && strings.HasPrefix(line, "auth "),
			settings.MSSFix > 0 && strings.HasPrefix(line, "mssfix "),
			!settings.IPv6 && hasPrefixOneOf(line, "tun-ipv6",
				`pull-filter ignore "route-ipv6"`,
				`pull-filter ignore "ifconfig-ipv6"`):
		default:
			modified = append(modified, line)
		}
	}

	// Add values
	modified = append(modified, connection.OpenVPNProtoLine())
	modified = append(modified, connection.OpenVPNRemoteLine())
	modified = append(modified, "dev "+settings.Interface)
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
		modified = append(modified, "persist-tun")
		modified = append(modified, "persist-key")
	}

	modified = append(modified, "") // trailing line

	return modified
}

func hasPrefixOneOf(s string, prefixes ...string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(s, prefix) {
			return true
		}
	}
	return false
}
