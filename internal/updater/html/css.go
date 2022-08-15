package html

import (
	"strings"

	"golang.org/x/net/html"
)

func HasClassStrings(node *html.Node, classStrings ...string) (match bool) {
	targetClasses := make(map[string]struct{}, len(classStrings))
	for _, classString := range classStrings {
		targetClasses[classString] = struct{}{}
	}

	classAttribute := Attribute(node, "class")
	classes := strings.Fields(classAttribute)
	for _, class := range classes {
		delete(targetClasses, class)
	}

	return len(targetClasses) == 0
}
