package html

import (
	"golang.org/x/net/html"
)

type MatchFunc func(node *html.Node) (match bool)

func MatchID(id string) MatchFunc {
	return func(node *html.Node) (match bool) {
		if node == nil {
			return false
		}

		return Attribute(node, "id") == id
	}
}

func MatchData(data string) MatchFunc {
	return func(node *html.Node) (match bool) {
		return node != nil && node.Type == html.ElementNode && node.Data == data
	}
}

func DirectChild(parent *html.Node,
	matchFunc MatchFunc) (child *html.Node) {
	for child := parent.FirstChild; child != nil; child = child.NextSibling {
		if matchFunc(child) {
			return child
		}
	}
	return nil
}

func DirectChildren(parent *html.Node,
	matchFunc MatchFunc) (children []*html.Node) {
	for child := parent.FirstChild; child != nil; child = child.NextSibling {
		if matchFunc(child) {
			children = append(children, child)
		}
	}
	return children
}
