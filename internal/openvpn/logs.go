package openvpn

import (
	"strings"

	"github.com/fatih/color"
	"github.com/qdm12/gluetun/internal/constants"
)

type logLevel uint8

const (
	levelInfo logLevel = iota
	levelWarn
	levelError
)

func processLogLine(s string) (filtered string, level logLevel) {
	for _, ignored := range []string{
		"WARNING: you are using user/group/chroot/setcon without persist-tun -- this may cause restarts to fail",
		"NOTE: UID/GID downgrade will be delayed because of --client, --pull, or --up-delay",
	} {
		if s == ignored {
			return "", levelInfo
		}
	}
	switch {
	case strings.HasPrefix(s, "NOTE: "):
		filtered = strings.TrimPrefix(s, "NOTE: ")
		level = levelInfo
	case strings.HasPrefix(s, "WARNING: "):
		filtered = strings.TrimPrefix(s, "WARNING: ")
		level = levelWarn
	case strings.HasPrefix(s, "ERROR: "):
		filtered = strings.TrimPrefix(s, "ERROR: ")
		level = levelError
	case strings.HasPrefix(s, "Options error: "):
		filtered = strings.TrimPrefix(s, "Options error: ")
		level = levelError
	case s == "Initialization Sequence Completed":
		return color.HiGreenString(s), levelInfo
	case s == "AUTH: Received control message: AUTH_FAILED":
		filtered = s + `

Your credentials might be wrong 🤨

`
		level = levelError
	case strings.Contains(s, "TLS Error: TLS key negotiation failed to occur within 60 seconds (check your network connectivity)"): //nolint:lll
		filtered = s + `
🚒🚒🚒🚒🚒🚨🚨🚨🚨🚨🚨🚒🚒🚒🚒🚒
That error usually happens because either:

1. The VPN server IP address you are trying to connect to is no longer valid 🔌
   Update your server information using https://github.com/qdm12/gluetun/wiki/Updating-Servers

2. The VPN server crashed 💥, try changing your VPN servers filtering options such as SERVER_REGIONS

3. Your Internet connection is not working 🤯, ensure it works

4. Something else ➡️ https://github.com/qdm12/gluetun/issues/new/choose
`
		level = levelWarn
	default:
		filtered = s
		level = levelInfo
	}

	switch {
	case filtered == "RTNETLINK answers: File exists":
		filtered = "OpenVPN tried to add an IP route which already exists (" + filtered + ")"
		level = levelWarn
	case strings.HasPrefix(filtered, "Linux route add command failed: "):
		filtered = "Previous error details: " + filtered
		level = levelWarn
	}

	filtered = constants.ColorOpenvpn().Sprintf(filtered)
	return filtered, level
}
