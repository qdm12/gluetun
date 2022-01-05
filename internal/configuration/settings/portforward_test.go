package settings

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_PortForwarding_String(t *testing.T) {
	t.Parallel()

	settings := PortForwarding{
		Enabled: boolPtr(false),
	}

	s := settings.String()

	assert.Empty(t, s)
}
