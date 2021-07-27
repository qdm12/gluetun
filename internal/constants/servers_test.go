package constants

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/qdm12/gluetun/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_parseAllServers(t *testing.T) {
	t.Parallel()

	servers, err := parseAllServers()

	require.NoError(t, err)
	require.NotEmpty(t, len(servers.Cyberghost.Servers))
}

func digestServerModelVersion(t *testing.T, server interface{}, version uint16) string {
	t.Helper()
	bytes, err := json.Marshal(server)
	if err != nil {
		t.Fatal(err)
	}
	bytes = append(bytes, []byte(fmt.Sprintf("%d", version))...)
	arr := sha256.Sum256(bytes)
	hexString := hex.EncodeToString(arr[:])
	if len(hexString) > 8 {
		hexString = hexString[:8]
	}
	return hexString
}

func Test_versions(t *testing.T) {
	t.Parallel()
	allServers := GetAllServers()
	const format = "you forgot to update the version for %s"
	testCases := map[string]struct {
		model   interface{}
		version uint16
		digest  string
	}{
		"Cyberghost": {
			model:   models.CyberghostServer{},
			version: allServers.Cyberghost.Version,
			digest:  "229828de",
		},
		"Fastestvpn": {
			model:   models.FastestvpnServer{},
			version: allServers.Fastestvpn.Version,
			digest:  "8825919b",
		},
		"HideMyAss": {
			model:   models.HideMyAssServer{},
			version: allServers.HideMyAss.Version,
			digest:  "a93b4057",
		},
		"Ipvanish": {
			model:   models.IpvanishServer{},
			version: allServers.Ipvanish.Version,
			digest:  "2eb80d28",
		},
		"Ivpn": {
			model:   models.IvpnServer{},
			version: allServers.Ivpn.Version,
			digest:  "2eb80d28",
		},
		"Mullvad": {
			model:   models.MullvadServer{},
			version: allServers.Mullvad.Version,
			digest:  "2a009192",
		},
		"Nordvpn": {
			model:   models.NordvpnServer{},
			version: allServers.Nordvpn.Version,
			digest:  "a3b5d609",
		},
		"Privado": {
			model:   models.PrivadoServer{},
			version: allServers.Privado.Version,
			digest:  "dba6736c",
		},
		"Private Internet Access": {
			model:   models.PIAServer{},
			version: allServers.Pia.Version,
			digest:  "91db9bc9",
		},
		"Privatevpn": {
			model:   models.PrivatevpnServer{},
			version: allServers.Privatevpn.Version,
			digest:  "cba13d78",
		},
		"Protonvpn": {
			model:   models.ProtonvpnServer{},
			version: allServers.Protonvpn.Version,
			digest:  "b964085b",
		},
		"Purevpn": {
			model:   models.PurevpnServer{},
			version: allServers.Purevpn.Version,
			digest:  "23f2d422",
		},
		"Surfshark": {
			model:   models.SurfsharkServer{},
			version: allServers.Surfshark.Version,
			digest:  "58de06d8",
		},
		"Torguard": {
			model:   models.TorguardServer{},
			version: allServers.Torguard.Version,
			digest:  "6eb9028e",
		},
		"VPN Unlimited": {
			model:   models.VPNUnlimitedServer{},
			version: allServers.VPNUnlimited.Version,
			digest:  "5cb51319",
		},
		"Vyprvpn": {
			model:   models.VyprvpnServer{},
			version: allServers.Vyprvpn.Version,
			digest:  "58de06d8",
		},
		"Windscribe": {
			model:   models.WindscribeServer{},
			version: allServers.Windscribe.Version,
			digest:  "6f6c16d6",
		},
	}
	for name, testCase := range testCases {
		name := name
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			digest := digestServerModelVersion(t, testCase.model, testCase.version)
			failureMessage := fmt.Sprintf(format, name)
			assert.Equal(t, testCase.digest, digest, failureMessage)
		})
	}
}
