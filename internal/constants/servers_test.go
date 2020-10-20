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
			digest:  "f1e01afe",
		},
		"Private Internet Access Old": {
			model:   models.PIAOldServer{},
			version: allServers.PiaOld.Version,
			digest:  "4e25ce4a",
		},
		"Purevpn": {
			model:   models.PurevpnServer{},
			version: allServers.Purevpn.Version,
			digest:  "cc1a2219",
		},
		"Surfshark": {
			model:   models.SurfsharkServer{},
			version: allServers.Surfshark.Version,
			digest:  "042bef64",
		},
		"Vyprvpn": {
			model:   models.VyprvpnServer{},
			version: allServers.Vyprvpn.Version,
			digest:  "042bef64",
		},
		"Windscribe": {
			model:   models.WindscribeServer{},
			version: allServers.Windscribe.Version,
			digest:  "042bef64",
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
			digest:    "160631de",
		},
		"Mullvad": {
			servers:   allServers.Mullvad.Servers,
			timestamp: allServers.Mullvad.Timestamp,
			digest:    "e1fee56f",
		},
		"Nordvpn": {
			servers:   allServers.Nordvpn.Servers,
			timestamp: allServers.Nordvpn.Timestamp,
			digest:    "9fc9a579",
		},
		"Private Internet Access": {
			servers:   allServers.Pia.Servers,
			timestamp: allServers.Pia.Timestamp,
			digest:    "1571e777",
		},
		"Private Internet Access Old": {
			servers:   allServers.PiaOld.Servers,
			timestamp: allServers.PiaOld.Timestamp,
			digest:    "3566a800",
		},
		"Purevpn": {
			servers:   allServers.Purevpn.Servers,
			timestamp: allServers.Purevpn.Timestamp,
			digest:    "cdf9b708",
		},
		"Surfshark": {
			servers:   allServers.Surfshark.Servers,
			timestamp: allServers.Surfshark.Timestamp,
			digest:    "79484811",
		},
		"Vyprvpn": {
			servers:   allServers.Vyprvpn.Servers,
			timestamp: allServers.Vyprvpn.Timestamp,
			digest:    "1992457c",
		},
		"Windscribe": {
			servers:   allServers.Windscribe.Servers,
			timestamp: allServers.Windscribe.Timestamp,
			digest:    "eacad593",
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
