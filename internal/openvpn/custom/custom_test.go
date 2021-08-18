package custom

import (
	"net"
	"os"
	"testing"

	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_BuildConfig(t *testing.T) {
	t.Parallel()

	file, err := os.CreateTemp("", "")
	require.NoError(t, err)
	defer removeFile(t, file.Name())
	defer file.Close()

	_, err = file.WriteString("remote 1.9.8.7\nkeep me\ncipher remove")
	require.NoError(t, err)

	err = file.Close()
	require.NoError(t, err)

	settings := configuration.OpenVPN{
		Cipher: "cipher",
		MSSFix: 999,
		Config: file.Name(),
	}

	lines, connection, err := BuildConfig(settings)
	assert.NoError(t, err)

	expectedLines := []string{
		"keep me",
		"proto udp",
		"remote 1.9.8.7 1194",
		"mute-replay-warnings",
		"auth-nocache",
		"pull-filter ignore \"auth-token\"",
		"auth-retry nointeract",
		"suppress-timestamps",
		"verb 0",
		"data-ciphers-fallback cipher",
		"data-ciphers cipher",
		"mssfix 999",
		"pull-filter ignore \"route-ipv6\"",
		"pull-filter ignore \"ifconfig-ipv6\"",
		"user ",
	}
	assert.Equal(t, expectedLines, lines)

	expectedConnection := models.OpenVPNConnection{
		IP:       net.IPv4(1, 9, 8, 7),
		Port:     1194,
		Protocol: constants.UDP,
	}
	assert.Equal(t, expectedConnection, connection)
}
