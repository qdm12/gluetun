package config

import (
	"testing"

	"github.com/qdm12/golibs/logging"
	"github.com/stretchr/testify/assert"
)

func Test_processLogLine(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		s        string
		filtered string
		level    logging.Level
	}{
		"empty string":  {"", "", logging.LevelInfo},
		"random string": {"asdasqdb", "asdasqdb", logging.LevelInfo},
		"openvpn unknown": {
			"message",
			"message",
			logging.LevelInfo},
		"openvpn note": {
			"NOTE: message",
			"message",
			logging.LevelInfo},
		"openvpn warning": {
			"WARNING: message",
			"message",
			logging.LevelWarn},
		"openvpn options error": {
			"Options error: message",
			"message",
			logging.LevelError},
		"openvpn ignored message": {
			"NOTE: UID/GID downgrade will be delayed because of --client, --pull, or --up-delay",
			"",
			logging.LevelDebug},
		"openvpn success": {
			"Initialization Sequence Completed",
			"Initialization Sequence Completed",
			logging.LevelInfo},
		"openvpn auth failed": {
			"AUTH: Received control message: AUTH_FAILED",
			"AUTH: Received control message: AUTH_FAILED\n\nYour credentials might be wrong 🤨\n\n",
			logging.LevelError},
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
