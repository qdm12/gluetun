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
