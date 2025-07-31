package settings

import (
	"os"
	"testing"
)

func TestWireguardCustomConfigFile(t *testing.T) {
	// Use the test config file from the repo, with bogus values only
	confPath := "../../wireguard/wg_test.conf"
	// Optionally, you could copy it to a temp file if needed for mutability
	wg := &Wireguard{
		CustomConfigFile: &confPath,
	}

	if wg.CustomConfigFile == nil || *wg.CustomConfigFile == "" {
		t.Error("CustomConfigFile should be set")
	}

	// Check the file exists and is readable
	if _, err := os.Stat(*wg.CustomConfigFile); err != nil {
		t.Errorf("Custom config file not found: %v", err)
	}

	// Optionally, parse the config and check for bogus values
	data, err := os.ReadFile(*wg.CustomConfigFile)
	if err != nil {
		t.Fatalf("failed to read config file: %v", err)
	}
	content := string(data)
	// Check for presence of bogus values (not real keys/endpoints)
	if !containsBogusWireguardValues(content) {
		t.Error("wg_test.conf should only contain bogus values for testing")
	}
}

// containsBogusWireguardValues checks for known bogus values in the config
func containsBogusWireguardValues(content string) bool {
	// These should match the bogus values in wg_test.conf
	bogusKeys := []string{
		"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=", // bogus private key
		"BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB=", // bogus public key
		"1.2.3.4/32",            // bogus address
		"8.8.8.8",               // bogus DNS
		"123.123.123.123:51820", // bogus endpoint
	}
	for _, v := range bogusKeys {
		if !contains(content, v) {
			return false
		}
	}
	return true
}

func contains(s, substr string) bool {
	return len(substr) == 0 || (len(s) >= len(substr) && (s == substr || contains(s[1:], substr) || contains(s[:len(s)-1], substr)))
}
