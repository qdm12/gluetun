package updater

import (
	"compress/gzip"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_extractDebURL(t *testing.T) {
	t.Parallel()

	file, err := os.Open("testdata/index.html.gz")
	require.NoError(t, err)
	defer file.Close()

	gzipReader, err := gzip.NewReader(file)
	require.NoError(t, err)
	defer gzipReader.Close()

	html, err := io.ReadAll(gzipReader)
	require.NoError(t, err)

	debURLs, err := extractDebURL(string(html), "https://www.purevpn.com/download/linux-vpn")
	require.NoError(t, err)

	const want = "https://dhnx3d2u57yhc.cloudfront.net/cross-platform/linux-gui/2.9.0/PureVPN_amd64.deb"

	assert.Equal(t, want, debURLs)
}
