package openvpn

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
		"empty string":  {"", "", logging.InfoLevel},
		"random string": {"asdasqdb", "asdasqdb", logging.InfoLevel},
		"openvpn unknown": {
			"message",
			"message",
			logging.InfoLevel},
		"openvpn note": {
			"NOTE: message",
			"message",
			logging.InfoLevel},
		"openvpn warning": {
			"WARNING: message",
			"message",
			logging.WarnLevel},
		"openvpn options error": {
			"Options error: message",
			"message",
			logging.ErrorLevel},
		"openvpn ignored message": {
			"NOTE: UID/GID downgrade will be delayed because of --client, --pull, or --up-delay",
			"",
			""},
		"openvpn success": {
			"Initialization Sequence Completed",
			"Initialization Sequence Completed",
			logging.InfoLevel},
		"openvpn auth failed": {
			"AUTH: Received control message: AUTH_FAILED",
			"AUTH: Received control message: AUTH_FAILED\n\nYour credentials might be wrong 🤨\n\n💡 If you use Private Internet Access, check https://github.com/qdm12/gluetun/issues/265\n\n", //nolint:lll
			logging.ErrorLevel},
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
