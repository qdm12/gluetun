package tun

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Tun(t *testing.T) {
	t.Parallel()

	path := getTempPath(t)

	tun := New()

	defer func() {
		err := os.RemoveAll(path)
		require.NoError(t, err)
	}()

	// No file check fail
	err := tun.Check(path)
	require.Error(t, err)
	expectedMessage := "TUN device is not available: open " + path + ": no such file or directory"
	require.Equal(t, expectedMessage, err.Error())

	// Create simple file
	file, err := os.Create(path)
	require.NoError(t, err)
	err = file.Close()
	require.NoError(t, err)

	// Simple file check fail
	err = tun.Check(path)
	require.Error(t, err)
	expectedMessage = "TUN file has an unexpected rdev: 0 instead of expected 2760"
	require.Equal(t, expectedMessage, err.Error())

	// Create TUN device fail as file exists
	err = tun.Create(path)
	require.Error(t, err)
	require.Equal(t, "cannot create TUN device file node: file exists", err.Error())

	// Remove simple file
	err = os.Remove(path)
	require.NoError(t, err)

	// Create TUN device success
	err = tun.Create(path)
	require.NoError(t, err)

	// Check TUN device success
	err = tun.Check(path)
	require.NoError(t, err)
}

func getTempPath(t *testing.T) (path string) {
	t.Helper()
	file, err := os.CreateTemp("", "")
	require.NoError(t, err)
	path = file.Name()
	err = file.Close()
	require.NoError(t, err)
	err = os.Remove(path)
	require.NoError(t, err)
	return path
}
