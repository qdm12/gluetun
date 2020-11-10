package logging

import (
	"testing"

	"github.com/qdm12/golibs/logging"
	"github.com/stretchr/testify/assert"
)

func Test_PostProcessLine(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		s        string
		filtered string
		level    logging.Level
	}{
		"empty string":  {"", "", logging.InfoLevel},
		"random string": {"asdasqdb", "asdasqdb", logging.InfoLevel},
		"unbound notice": {
			"unbound: [1594595249] unbound[75:0] notice: init module 0: validator",
			"unbound: init module 0: validator",
			logging.InfoLevel},
		"unbound info": {
			"unbound: [1594595249] unbound[75:0] info: init module 0: validator",
			"unbound: init module 0: validator",
			logging.InfoLevel},
		"unbound warn": {
			"unbound: [1594595249] unbound[75:0] warn: init module 0: validator",
			"unbound: init module 0: validator",
			logging.WarnLevel},
		"unbound error": {
			"unbound: [1594595249] unbound[75:0] error: init module 0: validator",
			"unbound: init module 0: validator",
			logging.ErrorLevel},
		"unbound unknown": {
			"unbound: [1594595249] unbound[75:0] BLA: init module 0: validator",
			"unbound: BLA: init module 0: validator",
			logging.ErrorLevel},
		"openvpn unknown": {
			"openvpn: message",
			"openvpn: message",
			logging.InfoLevel},
		"openvpn note": {
			"openvpn: NOTE: message",
			"openvpn: message",
			logging.InfoLevel},
		"openvpn warning": {
			"openvpn: WARNING: message",
			"openvpn: message",
			logging.WarnLevel},
		"openvpn options error": {
			"openvpn: Options error: message",
			"openvpn: message",
			logging.ErrorLevel},
		"openvpn ignored message": {
			"openvpn: NOTE: UID/GID downgrade will be delayed because of --client, --pull, or --up-delay",
			"",
			""},
		"openvpn success": {
			"openvpn: Initialization Sequence Completed",
			"openvpn: Initialization Sequence Completed",
			logging.InfoLevel},
		"openvpn auth failed": {
			"openvpn: AUTH: Received control message: AUTH_FAILED",
			"openvpn: AUTH: Received control message: AUTH_FAILED\n\n  (IF YOU ARE USING PIA servers, MAYBE CHECK OUT https://github.com/qdm12/gluetun/issues/265)\n", //nolint:lll
			logging.ErrorLevel},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			filtered, level := PostProcessLine(tc.s)
			assert.Equal(t, tc.filtered, filtered)
			assert.Equal(t, tc.level, level)
		})
	}
}
