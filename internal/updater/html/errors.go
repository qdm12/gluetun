package html

import (
	"bytes"
	"fmt"

	"golang.org/x/net/html"
)

func WrapError(sentinelError error, node *html.Node) error {
	return fmt.Errorf("%w: in HTML code: %s",
		sentinelError, mustRenderHTML(node))
}

func WrapWarning(warning string, node *html.Node) string {
	return fmt.Sprintf("%s: in HTML code: %s",
		warning, mustRenderHTML(node))
}

func mustRenderHTML(node *html.Node) (rendered string) {
	stringBuffer := bytes.NewBufferString("")
	err := html.Render(stringBuffer, node)
	if err != nil {
		panic(err)
	}
	return stringBuffer.String()
}
