package htmlutils

import (
	"container/list"
	"errors"
	"fmt"

	"golang.org/x/net/html"
)

var (
	ErrAttrNotFound = errors.New("matching attribute not found")
)

func GetText(n *html.Node) string {
	return n.FirstChild.Data
}

func GetAttr(n *html.Node, key string) (string, error) {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val, nil
		}
	}

	return "", ErrAttrNotFound
}

func CheckAttrMatch(n *html.Node, attrKey string, checkValue string) bool {
	attrValue, err := GetAttr(n, attrKey)
	return err == nil && attrValue == checkValue
}

func CheckID(n *html.Node, idValue string) bool {
	return CheckAttrMatch(n, "id", idValue)
}

func CheckNodeType(n *html.Node, tagType string) bool {
	return n.Type == html.ElementNode && n.Data == tagType
}

func GetFirstNodeByID(n *html.Node, idValue string) *html.Node {
	return bfs(n, func(n *html.Node) bool {
		return CheckID(n, idValue)
	})
}

func GetFirstNodeByType(n *html.Node, nodeType string) *html.Node {
	return bfs(n, func(n *html.Node) bool {
		return CheckNodeType(n, nodeType)
	})
}

func GetNodesByType(n *html.Node, nodeType string) []*html.Node {
	nodes := []*html.Node{}
	for childNode := n.FirstChild; childNode != nil; childNode = childNode.NextSibling {
		if CheckNodeType(childNode, nodeType) {
			nodes = append(nodes, childNode)
		}
	}
	return nodes
}

// branching first search: returns the node matching the match function
// and nil if no node is found.
func bfs(rootNode *html.Node,
	match func(node *html.Node) bool) (node *html.Node) {
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
