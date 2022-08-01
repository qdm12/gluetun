package updater

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"

	htmlutils "github.com/qdm12/gluetun/internal/updater/html"
	"golang.org/x/net/html"
)

var ErrHTTPStatusCode = errors.New("HTTP status code is not OK")

func fetchAndParseWebsite(ctx context.Context, client *http.Client) (
	hostToData map[string]serverData, err error) {
	const url = "https://www.slickvpn.com/locations/"
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create HTTP request: %w", err)
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("do HTTP request: %w", err)
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %d %s", ErrHTTPStatusCode, response.StatusCode, response.Status)
	}

	rootNode, err := html.Parse(response.Body)
	if err != nil {
		_ = response.Body.Close()
		return nil, fmt.Errorf("parsing HTML code: %w", err)
	}

	hostToData, err = parseHTML(rootNode)
	if err != nil {
		_ = response.Body.Close()
		return nil, fmt.Errorf("parsing HTML: %w", err)
	}

	err = response.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("closing response body: %w", err)
	}

	return hostToData, nil
}

type serverData struct {
	ovpnURL string
	country string
	region  string
	city    string
}

func parseHTML(rootNode *html.Node) (hostToData map[string]serverData, err error) {
	locationTableNode := htmlutils.GetFirstNodeByID(rootNode, "location-table")
	if locationTableNode == nil {
		return nil, fmt.Errorf("unable to find html node with matching id")
	}

	tBodyNode := htmlutils.GetFirstNodeByType(locationTableNode, "tbody")
	if tBodyNode == nil {
		return nil, fmt.Errorf("unable to find tbody tag inside location table")
	}

	rowNodes := htmlutils.GetNodesByType(tBodyNode, "tr")
	hostToData = make(map[string]serverData, len(rowNodes))

	for _, rowNode := range rowNodes {
		var hostname string
		var data serverData
		columnIndex := 0
		const (
			columnIndexContinent = 0
			columnIndexCountry   = 1
			columnIndexCity      = 2
			columnIndexConfig    = 3
		)
		for cellNode := rowNode.FirstChild; cellNode != nil; cellNode = cellNode.NextSibling {
			switch columnIndex {
			case columnIndexContinent:
				data.region = cellNode.FirstChild.Data
			case columnIndexCountry:
				data.country = cellNode.FirstChild.Data
			case columnIndexCity:
				data.city = cellNode.FirstChild.Data
			case columnIndexConfig:
				linkNodes := htmlutils.GetNodesByType(cellNode, "a")
				for _, linkNode := range linkNodes {
					if htmlutils.GetText(linkNode) != "OpenVPN" {
						continue
					}

					data.ovpnURL, err = htmlutils.GetAttr(linkNode, "href")
					if err != nil {
						return nil, fmt.Errorf("get attribute value: %w", err)
					}

					hostname, err = extractHostnameFromURL(data.ovpnURL)
					if err != nil {
						return nil, fmt.Errorf("extract hostname from url: %w", err)
					}

					break
				}
			}

			columnIndex++
			if columnIndex == columnIndexConfig+1 {
				break
			}
		}

		hostToData[hostname] = data
	}

	return hostToData, nil
}

var serverNameRegex = regexp.MustCompile(`^.+\/(?P<serverName>.+)\.ovpn$`)

var (
	ErrExtractHostnameFromURL = errors.New("cannot extract hostname from url")
)

func extractHostnameFromURL(url string) (hostname string, err error) {
	matches := serverNameRegex.FindStringSubmatch(url)
	const minMatches = 2
	if len(matches) < minMatches {
		return "", fmt.Errorf("%w: from %s", ErrExtractHostnameFromURL, url)
	}
	hostname = matches[1]
	return hostname, nil
}

var (
	ErrHTMLAttributeNotFound = errors.New("html attribute not found")
)

func getAttributeValue(node *html.Node, attributeKey string) (
	attributeValue string, err error) {
	for _, attr := range node.Attr {
		if attr.Key == attributeKey {
			return attr.Val, nil
		}
	}

	return "", fmt.Errorf("%w: for key %s", ErrHTMLAttributeNotFound, attributeKey)
}
