package html

import "golang.org/x/net/html"

func Attribute(node *html.Node, key string) (value string) {
	for _, attribute := range node.Attr {
		if attribute.Key == key {
			return attribute.Val
		}
	}
	return ""
}
