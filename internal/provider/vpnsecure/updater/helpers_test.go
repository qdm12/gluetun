package updater

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/net/html"
)

func parseTestHTML(t *testing.T, htmlString string) *html.Node {
	t.Helper()
	rootNode, err := html.Parse(strings.NewReader(htmlString))
	require.NoError(t, err)
	return rootNode
}

func parseTestDataIndexHTML(t *testing.T) *html.Node {
	t.Helper()

	data, err := os.ReadFile("testdata/index.html")
	require.NoError(t, err)

	return parseTestHTML(t, string(data))
}
