package logging

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/private-internet-access-docker/internal/constants"
)

var regularExpressions = struct { //nolint:gochecknoglobals
	unboundPrefix          *regexp.Regexp
	shadowsocksPrefix      *regexp.Regexp
	shadowsocksErrorPrefix *regexp.Regexp
	tinyproxyLoglevel      *regexp.Regexp
	tinyproxyPrefix        *regexp.Regexp
}{
	unboundPrefix:          regexp.MustCompile(`unbound: \[[0-9]{10}\] unbound\[[0-9]+:0\] `),
	shadowsocksPrefix:      regexp.MustCompile(`shadowsocks:[ ]+2[0-9]{3}\-[0-1][0-9]\-[0-3][0-9] [0-2][0-9]:[0-5][0-9]:[0-5][0-9] `),
	shadowsocksErrorPrefix: regexp.MustCompile(`shadowsocks error:[ ]+2[0-9]{3}\-[0-1][0-9]\-[0-3][0-9] [0-2][0-9]:[0-5][0-9]:[0-5][0-9] `),
	tinyproxyLoglevel:      regexp.MustCompile(`INFO|CONNECT|NOTICE|WARNING|ERROR|CRITICAL`),
	tinyproxyPrefix:        regexp.MustCompile(`tinyproxy: .+[ ]+(Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec) [0-3][0-9] [0-2][0-9]:[0-5][0-9]:[0-5][0-9] \[[0-9]+\]: `),
}

func PostProcessLine(s string) (filtered string, level logging.Level) {
	switch {
	case strings.HasPrefix(s, "openvpn: "):
		filtered = constants.ColorOpenvpn().Sprintf(s)
		return filtered, logging.InfoLevel
	case strings.HasPrefix(s, "unbound: "):
		prefix := regularExpressions.unboundPrefix.FindString(s)
		filtered = s[len(prefix):]
		switch {
		case strings.HasPrefix(filtered, "notice: "):
			filtered = strings.TrimPrefix(filtered, "notice: ")
			level = logging.InfoLevel
		case strings.HasPrefix(filtered, "info: "):
			filtered = strings.TrimPrefix(filtered, "info: ")
			level = logging.InfoLevel
		case strings.HasPrefix(filtered, "warn: "):
			filtered = strings.TrimPrefix(filtered, "warn: ")
			level = logging.WarnLevel
		case strings.HasPrefix(filtered, "error: "):
			filtered = strings.TrimPrefix(filtered, "error: ")
			level = logging.ErrorLevel
		default:
			level = logging.ErrorLevel
		}
		filtered = fmt.Sprintf("unbound: %s", filtered)
		filtered = constants.ColorUnbound().Sprintf(filtered)
		return filtered, level
	case strings.HasPrefix(s, "shadowsocks: "):
		prefix := regularExpressions.shadowsocksPrefix.FindString(s)
		filtered = s[len(prefix):]
		switch {
		case strings.HasPrefix(filtered, "INFO: "):
			level = logging.InfoLevel
			filtered = strings.TrimPrefix(filtered, "INFO: ")
		default:
			level = logging.WarnLevel
		}
		filtered = fmt.Sprintf("shadowsocks: %s", filtered)
		filtered = constants.ColorShadowsocks().Sprintf(filtered)
		return filtered, level
	case strings.HasPrefix(s, "shadowsocks error: "):
		if strings.Contains(s, "ERROR: unable to resolve") { // caused by DNS blocking
			return "", logging.ErrorLevel
		}
		prefix := regularExpressions.shadowsocksErrorPrefix.FindString(s)
		filtered = s[len(prefix):]
		filtered = strings.TrimPrefix(filtered, "ERROR: ")
		filtered = fmt.Sprintf("shadowsocks: %s", filtered)
		filtered = constants.ColorShadowsocksError().Sprintf(filtered)
		return filtered, logging.ErrorLevel
	case strings.HasPrefix(s, "tinyproxy: "):
		logLevel := regularExpressions.tinyproxyLoglevel.FindString(s)
		prefix := regularExpressions.tinyproxyPrefix.FindString(s)
		filtered = fmt.Sprintf("tinyproxy: %s", s[len(prefix):])
		filtered = constants.ColorTinyproxy().Sprintf(filtered)
		switch logLevel {
		case "INFO", "CONNECT", "NOTICE":
			return filtered, logging.InfoLevel
		case "WARNING":
			return filtered, logging.WarnLevel
		case "ERROR", "CRITICAL":
			return filtered, logging.ErrorLevel
		default:
			return filtered, logging.ErrorLevel
		}
	}
	return s, logging.InfoLevel
}
