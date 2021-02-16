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
			digest:  "fd6242bb",
		},
		"Mullvad": {
			model:   models.MullvadServer{},
			version: allServers.Mullvad.Version,
			digest:  "665e9dc1",
		},
		"Nordvpn": {
			model:   models.NordvpnServer{},
			version: allServers.Nordvpn.Version,
			digest:  "040de8d0",
		},
		"Private Internet Access": {
			model:   models.PIAServer{},
			version: allServers.Pia.Version,
			digest:  "b90147aa",
		},
		"Privado": {
			model:   models.PrivadoServer{},
			version: allServers.Privado.Version,
			digest:  "1d5aeb23",
		},
		"Purevpn": {
			model:   models.PurevpnServer{},
			version: allServers.Purevpn.Version,
			digest:  "ada45379",
		},
		"Surfshark": {
			model:   models.SurfsharkServer{},
			version: allServers.Surfshark.Version,
			digest:  "042bef64",
		},
		"Torguard": {
			model:   models.TorguardServer{},
			version: allServers.Torguard.Version,
			digest:  "752702f3",
		},
		"Vyprvpn": {
			model:   models.VyprvpnServer{},
			version: allServers.Vyprvpn.Version,
			digest:  "042bef64",
		},
		"Windscribe": {
			model:   models.WindscribeServer{},
			version: allServers.Windscribe.Version,
			digest:  "6e3ca639",
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
			digest:    "5d3a8cbf",
		},
		"Mullvad": {
			servers:   allServers.Mullvad.Servers,
			timestamp: allServers.Mullvad.Timestamp,
			digest:    "e2e006cf",
		},
		"Nordvpn": {
			servers:   allServers.Nordvpn.Servers,
			timestamp: allServers.Nordvpn.Timestamp,
			digest:    "2296312c",
		},
		"Private Internet Access": {
			servers:   allServers.Pia.Servers,
			timestamp: allServers.Pia.Timestamp,
			digest:    "1d2938a1",
		},
		"Purevpn": {
			servers:   allServers.Purevpn.Servers,
			timestamp: allServers.Purevpn.Timestamp,
			digest:    "cd19edf5",
		},
		"Privado": {
			servers:   allServers.Privado.Servers,
			timestamp: allServers.Privado.Timestamp,
			digest:    "2ac55360",
		},
		"Surfshark": {
			servers:   allServers.Surfshark.Servers,
			timestamp: allServers.Surfshark.Timestamp,
			digest:    "1a7f38bb",
		},
		"Torguard": {
			servers:   allServers.Torguard.Servers,
			timestamp: allServers.Torguard.Timestamp,
			digest:    "dffab93e",
		},
		"Vyprvpn": {
			servers:   allServers.Vyprvpn.Servers,
			timestamp: allServers.Vyprvpn.Timestamp,
			digest:    "1753d4f8",
		},
		"Windscribe": {
			servers:   allServers.Windscribe.Servers,
			timestamp: allServers.Windscribe.Timestamp,
			digest:    "4e719aa3",
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
