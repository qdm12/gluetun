package updater

import (
	"bytes"
	"container/list"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/common"
	"golang.org/x/net/html"
)

func fetchServers(ctx context.Context, client *http.Client,
	warner common.Warner) (servers []models.Server, err error) {
	rootNode, err := fetchHTML(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("fetching HTML code: %w", err)
	}

	servers, warnings, err := parseHTML(rootNode)
	for _, warning := range warnings {
		warner.Warn(warning)
	}
	if err != nil {
		return nil, fmt.Errorf("parsing HTML code: %w", err)
	}

	return servers, nil
}

var ErrHTTPStatusCode = errors.New("HTTP status code is not OK")

func fetchHTML(ctx context.Context, client *http.Client) (rootNode *html.Node, err error) {
	const url = "https://www.vpnsecure.me/vpn-locations/"
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %d %s",
			ErrHTTPStatusCode, response.StatusCode, response.Status)
	}

	rootNode, err = html.Parse(response.Body)
	if err != nil {
		_ = response.Body.Close()
		return nil, fmt.Errorf("parsing HTML code: %w", err)
	}

	err = response.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("closing response body: %w", err)
	}

	return rootNode, nil
}

var (
	ErrHTMLServersDivNotFound = errors.New("HTML servers container div not found")
)

const divString = "div"

func parseHTML(rootNode *html.Node) (servers []models.Server,
	warnings []string, err error) {
	// Find div container for all servers, searching with BFS.
	serversDiv := findServersDiv(rootNode)
	if serversDiv == nil {
		return nil, nil, wrapHTMLError(ErrHTMLServersDivNotFound, rootNode)
	}

	for countryNode := serversDiv.FirstChild; countryNode != nil; countryNode = countryNode.NextSibling {
		if countryNode.Data != divString {
			// empty line(s) and tab(s)
			continue
		}

		country := findCountry(countryNode)
		if country == "" {
			warnings = append(warnings, wrapHTMLWarning("country not found", countryNode))
			continue
		}

		grid := bfs(countryNode, matchGridDiv)
		if grid == nil {
			warnings = append(warnings, wrapHTMLWarning("grid div not found", countryNode))
			continue
		}

		gridItems := getDirectChildren(grid, matchGridItem)
		if len(gridItems) == 0 {
			warnings = append(warnings, wrapHTMLWarning("no grid item found", grid))
			continue
		}

		for _, gridItem := range gridItems {
			server, warning := parseHTMLGridItem(gridItem)
			if warning != "" {
				warnings = append(warnings, warning)
				continue
			}

			server.Country = country
			servers = append(servers, server)
		}
	}

	return servers, warnings, nil
}

func parseHTMLGridItem(gridItem *html.Node) (
	server models.Server, warning string) {
	gridItemDT := getDirectChild(gridItem, matchDT)
	if gridItemDT == nil {
		return server, wrapHTMLWarning("grid item <dt> not found", gridItem)
	}

	host := findHost(gridItemDT)
	if host == "" {
		return server, wrapHTMLWarning("host not found", gridItemDT)
	}

	status := findStatus(gridItemDT)
	if !strings.EqualFold(status, "up") {
		warning := fmt.Sprintf("skipping server with host %s which has status %q", host, status)
		warning = wrapHTMLWarning(warning, gridItemDT)
		return server, warning
	}

	gridItemDD := getDirectChild(gridItem, matchDD)
	if gridItemDD == nil {
		return server, wrapHTMLWarning("grid item dd not found", gridItem)
	}

	region := findSpanStrong(gridItemDD, "Region:")
	if region == "" {
		warning := fmt.Sprintf("region for host %s not found", host)
		return server, wrapHTMLWarning(warning, gridItemDD)
	}

	city := findSpanStrong(gridItemDD, "City:")
	if city == "" {
		warning := fmt.Sprintf("region for host %s not found", host)
		return server, wrapHTMLWarning(warning, gridItemDD)
	}

	premiumString := findSpanStrong(gridItemDD, "Premium:")
	if premiumString == "" {
		warning := fmt.Sprintf("premium for host %s not found", host)
		return server, wrapHTMLWarning(warning, gridItemDD)
	}

	return models.Server{
		Region:   region,
		City:     city,
		Hostname: host + ".isponeder.com",
		Premium:  strings.EqualFold(premiumString, "yes"),
	}, ""
}

func findCountry(countryNode *html.Node) (country string) {
	for node := countryNode.FirstChild; node != nil; node = node.NextSibling {
		if node.Data != "a" {
			continue
		}
		for subNode := node.FirstChild; subNode != nil; subNode = subNode.NextSibling {
			if subNode.Data != "h4" {
				continue
			}
			return subNode.FirstChild.Data
		}
	}
	return ""
}

type matchFunc func(node *html.Node) (match bool)

func findServersDiv(rootNode *html.Node) (serversDiv *html.Node) {
	locationsDiv := bfs(rootNode, matchLocationsListDiv)
	if locationsDiv == nil {
		return nil
	}

	return bfs(locationsDiv, matchServersDiv)
}

func findHost(gridItemDT *html.Node) (host string) {
	for child := gridItemDT.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.TextNode {
			data := strings.TrimSpace(child.Data)
			if data == "" {
				continue // empty lines
			}
			return data
		}
	}
	return ""
}

func findStatus(gridItemDT *html.Node) (status string) {
	for child := gridItemDT.FirstChild; child != nil; child = child.NextSibling {
		if child.Data == "span" && containsClassString(child, "status") {
			return strings.TrimSpace(child.FirstChild.Data)
		}
	}
	return ""
}

func matchServersDiv(node *html.Node) (match bool) {
	return node != nil && node.Data == divString &&
		containsClassString(node, "blk__i")
}

func matchLocationsListDiv(node *html.Node) (match bool) {
	return node != nil && node.Data == divString &&
		containsClassString(node, "locations-list")
}

func matchGridDiv(node *html.Node) (match bool) {
	return node != nil && node.Data == divString &&
		containsClassString(node, "grid--locations")
}

func matchGridItem(node *html.Node) (match bool) {
	return node != nil && node.Data == "dl" &&
		containsClassString(node, "grid__i")
}

func matchDT(node *html.Node) (match bool) {
	return node != nil && node.Data == "dt"
}

func matchDD(node *html.Node) (match bool) {
	return node != nil && node.Data == "dd"
}

func findSpanStrong(gridItemDD *html.Node, spanData string) (
	strongValue string) {
	spanFound := false
	for child := gridItemDD.FirstChild; child != nil; child = child.NextSibling {
		if !matchDiv(child) {
			continue
		}

		for subchild := child.FirstChild; subchild != nil; subchild = subchild.NextSibling {
			if matchSpan(subchild) && subchild.FirstChild.Data == spanData {
				spanFound = true
				break
			}
		}

		if !spanFound {
			continue
		}

		for subchild := child.FirstChild; subchild != nil; subchild = subchild.NextSibling {
			if matchStrong(subchild) {
				return subchild.FirstChild.Data
			}
		}
	}

	return ""
}

func matchDiv(node *html.Node) (match bool) {
	return node != nil && node.Data == "div"
}

func matchSpan(node *html.Node) (match bool) {
	return node != nil && node.Data == "span"
}

func matchStrong(node *html.Node) (match bool) {
	return node != nil && node.Data == "strong"
}

func getDirectChild(parent *html.Node,
	matchFunc matchFunc) (child *html.Node) {
	for child := parent.FirstChild; child != nil; child = child.NextSibling {
		if matchFunc(child) {
			return child
		}
	}
	return nil
}

func getDirectChildren(parent *html.Node,
	matchFunc matchFunc) (children []*html.Node) {
	for child := parent.FirstChild; child != nil; child = child.NextSibling {
		if matchFunc(child) {
			children = append(children, child)
		}
	}
	return children
}

func containsClassString(node *html.Node, classString string) (match bool) {
	attributes := node.Attr
	for _, attribute := range attributes {
		if attribute.Key != "class" {
			continue
		}

		if strings.Contains(attribute.Val, classString) {
			return true
		}
	}

	return false
}

// bfs returns the node matching the match function and nil
// if no node is found.
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

func wrapHTMLError(sentinelError error, node *html.Node) error {
	return fmt.Errorf("%w: in HTML code: %s",
		sentinelError, mustRenderHTML(node))
}

func wrapHTMLWarning(warning string, node *html.Node) string {
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
