package openvpn

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_processLogLine(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		s        string
		filtered string
		level    logLevel
	}{
		"empty string":  {"", "", levelInfo},
		"random string": {"asdasqdb", "asdasqdb", levelInfo},
		"openvpn unknown": {
			"message",
			"message",
			levelInfo},
		"openvpn note": {
			"NOTE: message",
			"message",
			levelInfo},
		"openvpn warning": {
			"WARNING: message",
			"message",
			levelWarn},
		"openvpn options error": {
			"Options error: message",
			"message",
			levelError},
		"openvpn ignored message": {
			"NOTE: UID/GID downgrade will be delayed because of --client, --pull, or --up-delay",
			"",
			levelInfo},
		"openvpn success": {
			"Initialization Sequence Completed",
			"Initialization Sequence Completed",
			levelInfo},
		"openvpn auth failed": {
			"AUTH: Received control message: AUTH_FAILED",
			"AUTH: Received control message: AUTH_FAILED\n\nYour credentials might be wrong 🤨\n\n",
			levelError},
		"TLS key negotiation error": {
			s: "TLS Error: TLS key negotiation failed to occur within " +
				"60 seconds (check your network connectivity)",
			filtered: "TLS Error: TLS key negotiation failed to occur within " +
				"60 seconds (check your network connectivity)" + `
🚒🚒🚒🚒🚒🚨🚨🚨🚨🚨🚨🚒🚒🚒🚒🚒
That error usually happens because either:

1. The VPN server IP address you are trying to connect to is no longer valid 🔌
   Update your server information using https://github.com/qdm12/gluetun/wiki/Updating-Servers

2. The VPN server crashed 💥, try changing your VPN servers filtering options such as SERVER_REGIONS

3. Your Internet connection is not working 🤯, ensure it works

4. Something else ➡️ https://github.com/qdm12/gluetun/issues/new/choose
`,
			level: levelWarn,
		},
		"RTNETLINK answers: File exists": {
			s: "ERROR: RTNETLINK answers: File exists",
			filtered: "OpenVPN tried to add an IP route which already exists " +
				"(RTNETLINK answers: File exists)",
			level: levelWarn,
		},
		"Linux route add command failed": {
			s:        "ERROR: Linux route add command failed: some error",
			filtered: "Previous error details: Linux route add command failed: some error",
			level:    levelWarn,
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			filtered, level := processLogLine(tc.s)
			assert.Equal(t, tc.filtered, filtered)
			assert.Equal(t, tc.level, level)
		})
	}
}
