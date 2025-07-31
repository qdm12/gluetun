package settings

import (
	"os"
	"strings"
	"testing"
)

func TestWireguardCustomConfigFile(t *testing.T) {
	// Use the test config file from the repo, with bogus values only
	confPath := "../../wireguard/wg_test.conf"
	t.Logf("Using config file: %s", confPath)
	wg := &Wireguard{
		CustomConfigFile: &confPath,
	}

	t.Log("Checking if CustomConfigFile is set...")
	if wg.CustomConfigFile == nil || *wg.CustomConfigFile == "" {
		t.Error("CustomConfigFile should be set")
	} else {
		t.Logf("CustomConfigFile is set to: %s", *wg.CustomConfigFile)
	}

	t.Log("Checking if config file exists and is readable...")
	if _, err := os.Stat(*wg.CustomConfigFile); err != nil {
		t.Errorf("Custom config file not found: %v", err)
	} else {
		t.Log("Config file exists and is readable.")
	}

	t.Log("Reading config file contents...")
	data, err := os.ReadFile(*wg.CustomConfigFile)
	if err != nil {
		t.Fatalf("failed to read config file: %v", err)
	}
	content := string(data)
	t.Logf("Config file content (first 100 chars): %.100s", content)

	t.Log("Checking for bogus or real values in config file...")
	if containsBogusWireguardValues(content) {
		t.Log("Expected bogus test values found in config file (safe for public use/testing).")
	} else {
		t.Error("wg_test.conf does not contain the expected bogus test values. " +
			"The test config should contain the known safe test values, not real production values.")
	}
}

// containsBogusWireguardValues checks for known bogus values in the config
func containsBogusWireguardValues(content string) bool {
	bogusValues := []string{
		"eOjmwWHYOpRwZHvFgxNCD/msMO+PGnF3xZejEei8Klw=", // bogus private key from wg_test.conf
		"g+04U6gzk1+3zNUSMQtRvILfToBT8r6gsR0hEb4BOWI=", // bogus public key from wg_test.conf
		"1.2.3.4/32",            // bogus address from wg_test.conf
		"8.8.8.8",               // bogus DNS from wg_test.conf
		"123.123.123.123:51820", // bogus endpoint from wg_test.conf
	}
	
	// Check if ALL bogus values are present (indicates test config)
	for _, bogusValue := range bogusValues {
		if !strings.Contains(content, bogusValue) {
			return false
		}
	}
	return true
}
