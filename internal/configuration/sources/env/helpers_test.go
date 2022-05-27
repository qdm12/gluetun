package env

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setTestEnv is used to set environment variables in
// parallel tests.
func setTestEnv(t *testing.T, key, value string) {
	t.Helper()
	existing := os.Getenv(key)
	err := os.Setenv(key, value) //nolint:tenv
	t.Cleanup(func() {
		err = os.Setenv(key, existing)
		assert.NoError(t, err)
	})
	require.NoError(t, err)
}
