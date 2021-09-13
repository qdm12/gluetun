package extract

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func removeFile(t *testing.T, filename string) {
	t.Helper()
	err := os.RemoveAll(filename)
	require.NoError(t, err)
}
