package settings

import (
	"os"
	"testing"
)

func TestWireguardCustomConfigFile(t *testing.T) {
	const conf = `[Interface]
PrivateKey = MKhPTJVe2DVEkOTdI6jSDXF6h25V+/mx+P/p1QF11lU=
Address = 10.2.0.2/32
DNS = 10.2.0.1

[Peer]
PublicKey = JobdIlwHN75a1pPOqfUNu0a1eQXcwqaz5vyrq0qT6Ek=
AllowedIPs = 0.0.0.0/0, ::/0
Endpoint = 45.83.137.1:51820
`

	file, err := os.CreateTemp("", "wgtest-*.conf")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(file.Name())
	_, err = file.WriteString(conf)
	if err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	file.Close()

	fileName := file.Name()
	wg := &Wireguard{
		CustomConfigFile: &fileName,
	}

	if wg.CustomConfigFile == nil || *wg.CustomConfigFile == "" {
		t.Error("CustomConfigFile should be set")
	}

	// Here you would call the logic that parses and uses the custom config file.
	// For now, just check the file exists and is readable.
	if _, err := os.Stat(*wg.CustomConfigFile); err != nil {
		t.Errorf("Custom config file not found: %v", err)
	}
}
