package html

import (
	"container/list"
	"fmt"

	"golang.org/x/net/html"
)

// BFS returns the node matching the match function and nil
// if no node is found.
func BFS(rootNode *html.Node, match MatchFunc) (node *html.Node) {
	visited := make(map[*html.Node]struct{})
	queue := list.New()
	_ = queue.PushBack(rootNode)

	for queue.Len() > 0 {
		listElement := queue.Front()
		node, ok := queue.Remove(listElement).(*html.Node)
		if !ok {
			panic(fmt.Sprintf("linked list has bad type %T", listElement.Value))
		}

		if node == nil {
			continue
		}

		if _, ok := visited[node]; ok {
			continue
		}
		visited[node] = struct{}{}

		if match(node) {
			return node
		}

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			_ = queue.PushBack(child)
		}
	}

	return nil
}
