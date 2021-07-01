package constants

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/qdm12/gluetun/internal/models"
	"github.com/stretchr/testify/assert"
)

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
			digest:  "0bd93da1",
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

func digestServersTimestamp(t *testing.T, servers interface{}, timestamp int64) string {
	t.Helper()
	bytes, err := json.Marshal(servers)
	if err != nil {
		t.Fatal(err)
	}
	bytes = append(bytes, []byte(fmt.Sprintf("%d", timestamp))...)
	arr := sha256.Sum256(bytes)
	hexString := hex.EncodeToString(arr[:])
	if len(hexString) > 8 {
		hexString = hexString[:8]
	}
	return hexString
}

func Test_timestamps(t *testing.T) {
	t.Parallel()
	allServers := GetAllServers()
	const format = "you forgot to update the timestamp for %s"
	testCases := map[string]struct {
		servers   interface{}
		timestamp int64
		digest    string
	}{
		"Cyberghost": {
			servers:   allServers.Cyberghost.Servers,
			timestamp: allServers.Cyberghost.Timestamp,
			digest:    "b3ca3118",
		},
		"Fastestvpn": {
			servers:   allServers.Fastestvpn.Version,
			timestamp: allServers.Fastestvpn.Timestamp,
			digest:    "f0ef6b0b",
		},
		"HideMyAss": {
			servers:   allServers.HideMyAss.Servers,
			timestamp: allServers.HideMyAss.Timestamp,
			digest:    "8f872ac4",
		},
		"Ipvanish": {
			servers:   allServers.Ipvanish.Servers,
			timestamp: allServers.Ipvanish.Timestamp,
			digest:    "7c60dc5d",
		},
		"Ivpn": {
			servers:   allServers.Ivpn.Servers,
			timestamp: allServers.Ivpn.Timestamp,
			digest:    "42f92754",
		},
		"Mullvad": {
			servers:   allServers.Mullvad.Servers,
			timestamp: allServers.Mullvad.Timestamp,
			digest:    "01f2315f",
		},
		"Nordvpn": {
			servers:   allServers.Nordvpn.Servers,
			timestamp: allServers.Nordvpn.Timestamp,
			digest:    "b2619eea",
		},
		"Privado": {
			servers:   allServers.Privado.Servers,
			timestamp: allServers.Privado.Timestamp,
			digest:    "df378478",
		},
		"Private Internet Access": {
			servers:   allServers.Pia.Servers,
			timestamp: allServers.Pia.Timestamp,
			digest:    "ad230e95",
		},
		"Privatevpn": {
			servers:   allServers.Privatevpn.Servers,
			timestamp: allServers.Privatevpn.Timestamp,
			digest:    "e8d8255a",
		},
		"Protonvpn": {
			servers:   allServers.Protonvpn.Servers,
			timestamp: allServers.Protonvpn.Timestamp,
			digest:    "3b86393e",
		},
		"Purevpn": {
			servers:   allServers.Purevpn.Servers,
			timestamp: allServers.Purevpn.Timestamp,
			digest:    "9263cfdb",
		},
		"Surfshark": {
			servers:   allServers.Surfshark.Servers,
			timestamp: allServers.Surfshark.Timestamp,
			digest:    "55669f49",
		},
		"Torguard": {
			servers:   allServers.Torguard.Servers,
			timestamp: allServers.Torguard.Timestamp,
			digest:    "af54b9e8",
		},
		"VPN Unlimited": {
			servers:   allServers.VPNUnlimited.Servers,
			timestamp: allServers.VPNUnlimited.Timestamp,
			digest:    "f6ddb84c",
		},
		"Vyprvpn": {
			servers:   allServers.Vyprvpn.Servers,
			timestamp: allServers.Vyprvpn.Timestamp,
			digest:    "a62c484b",
		},
		"Windscribe": {
			servers:   allServers.Windscribe.Servers,
			timestamp: allServers.Windscribe.Timestamp,
			digest:    "d9290941",
		},
	}
	for name, testCase := range testCases {
		name := name
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			digest := digestServersTimestamp(t, testCase.servers, testCase.timestamp)
			failureMessage := fmt.Sprintf(format, name)
			assert.Equal(t, testCase.digest, digest, failureMessage)
		})
	}
}
