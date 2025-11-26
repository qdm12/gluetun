package internal

import (
	"context"
	"fmt"
)

func MullvadTest(ctx context.Context, logger Logger) error {
	expectedSecrets := []string{
		"Wireguard private key",
		"Wireguard address",
	}
	secrets, err := readSecrets(ctx, expectedSecrets, logger)
	if err != nil {
		return fmt.Errorf("reading secrets: %w", err)
	}

	env := []string{
		"VPN_SERVICE_PROVIDER=mullvad",
		"VPN_TYPE=wireguard",
		"LOG_LEVEL=debug",
		"SERVER_COUNTRIES=USA",
		"WIREGUARD_PRIVATE_KEY=" + secrets[0],
		"WIREGUARD_ADDRESSES=" + secrets[1],
	}
	return simpleTest(ctx, env, logger)
}
