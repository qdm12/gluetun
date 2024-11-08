package netlink

import (
	"testing"

	"github.com/qdm12/log"
)

func Test_IsIPv6Supported(t *testing.T) {
	n := New(log.New(log.SetLevel(log.LevelDebug)))
	supported, err := n.IsIPv6Supported()
	t.Log(supported, err)
}
