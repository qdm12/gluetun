package updater

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/andybalholm/cascadia"
	"golang.org/x/net/html"
)

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

	hostToData, err = parseHTML(response.Body)
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

var (
	locationTableSelector = cascadia.MustCompile(`table#location-table > tbody > tr`) //nolint:gochecknoglobals
	ovpnHrefSelector      = cascadia.MustCompile(`a[href$='.ovpn']`)                  //nolint:gochecknoglobals
)

func parseHTML(body io.Reader) (hostToData map[string]serverData, err error) {
	root, err := html.Parse(body)
	if err != nil {
		return nil, fmt.Errorf("parsing response body: %w", err)
	}

	rowNodes := cascadia.QueryAll(root, locationTableSelector)
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
				// TODO Translate continent to region
				data.region = cellNode.FirstChild.Data
			case columnIndexCountry:
				data.country = cellNode.FirstChild.Data
			case columnIndexCity:
				data.city = cellNode.FirstChild.Data
			case columnIndexConfig:
				linkNodes := cascadia.QueryAll(cellNode, ovpnHrefSelector)
				for _, linkNode := range linkNodes {
					if !strings.EqualFold(linkNode.FirstChild.Data, "openvpn") {
						continue
					}

					data.ovpnURL, err = getAttributeValue(linkNode, "href")
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
