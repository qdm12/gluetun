package extract

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_readCustomConfigLines(t *testing.T) {
	t.Parallel()

	file, err := os.CreateTemp("", "")
	require.NoError(t, err)
	defer removeFile(t, file.Name())
	defer file.Close()

	_, err = file.WriteString("line one\nline two\nline three\n")
	require.NoError(t, err)

	err = file.Close()
	require.NoError(t, err)

	lines, err := readCustomConfigLines(file.Name())
	assert.NoError(t, err)

	expectedLines := []string{
		"line one", "line two", "line three", "",
	}
	assert.Equal(t, expectedLines, lines)
}
