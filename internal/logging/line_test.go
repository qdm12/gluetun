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
		"shadowsocks stdout info": {
			"shadowsocks:  2020-07-12 23:07:25 INFO: UDP relay enabled",
			"shadowsocks: UDP relay enabled",
			logging.InfoLevel},
		"shadowsocks stdout other": {
			"shadowsocks:  2020-07-12 23:07:25 BLABLA: UDP relay enabled",
			"shadowsocks: BLABLA: UDP relay enabled",
			logging.WarnLevel},
		"shadowsocks stderr": {
			"shadowsocks error:  2020-07-12 23:07:25 Some error",
			"shadowsocks: Some error",
			logging.ErrorLevel},
		"shadowsocks stderr unable to resolve muted": {
			"shadowsocks error:  2020-07-12 23:07:25 ERROR: unable to resolve",
			"",
			logging.ErrorLevel},
		"tinyproxy info": {
			"tinyproxy: INFO      Jul 12 23:07:25 [32]: Reloading config file",
			"tinyproxy: Reloading config file",
			logging.InfoLevel},
		"tinyproxy connect": {
			"tinyproxy: CONNECT      Jul 12 23:07:25 [32]: Reloading config file",
			"tinyproxy: Reloading config file",
			logging.InfoLevel},
		"tinyproxy notice": {
			"tinyproxy: NOTICE      Jul 12 23:07:25 [32]: Reloading config file",
			"tinyproxy: Reloading config file",
			logging.InfoLevel},
		"tinyproxy warning": {
			"tinyproxy: WARNING      Jul 12 23:07:25 [32]: Reloading config file",
			"tinyproxy: Reloading config file",
			logging.WarnLevel},
		"tinyproxy error": {
			"tinyproxy: ERROR      Jul 12 23:07:25 [32]: Reloading config file",
			"tinyproxy: Reloading config file",
			logging.ErrorLevel},
		"tinyproxy critical": {
			"tinyproxy: CRITICAL      Jul 12 23:07:25 [32]: Reloading config file",
			"tinyproxy: Reloading config file",
			logging.ErrorLevel},
		"tinyproxy unknown": {
			"tinyproxy: BLABLA      Jul 12 23:07:25 [32]: Reloading config file",
			"tinyproxy: Reloading config file",
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
