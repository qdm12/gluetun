package dns

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
		"unbound notice": {
			"[1594595249] unbound[75:0] notice: init module 0: validator",
			"init module 0: validator",
			levelInfo},
		"unbound info": {
			"[1594595249] unbound[75:0] info: init module 0: validator",
			"init module 0: validator",
			levelInfo},
		"unbound warn": {
			"[1594595249] unbound[75:0] warn: init module 0: validator",
			"init module 0: validator",
			levelWarn},
		"unbound error": {
			"[1594595249] unbound[75:0] error: init module 0: validator",
			"init module 0: validator",
			levelError},
		"unbound unknown": {
			"[1594595249] unbound[75:0] BLA: init module 0: validator",
			"BLA: init module 0: validator",
			levelInfo},
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
