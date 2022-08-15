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

func fetchServers(ctx context.Context, client *http.Client) (
	hostToData map[string]serverData, err error) {
	const url = "https://www.slickvpn.com/locations/"
	rootNode, err := htmlutils.Fetch(ctx, client, url)
	if err != nil {
		return nil, fmt.Errorf("fetching HTML code: %w", err)
	}

	hostToData, err = parseHTML(rootNode)
	if err != nil {
		return nil, fmt.Errorf("parsing HTML: %w", err)
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
	ErrLocationTableNotFound = errors.New("HTML location table node not found")
	ErrTbodyNotFound         = errors.New("HTML tbody node not found")
	ErrExtractOpenVPNURL     = errors.New("failed extracting OpenVPN URL")
)

func parseHTML(rootNode *html.Node) (hostToData map[string]serverData, err error) {
	locationTableNode := htmlutils.BFS(rootNode, matchLocationTable)
	if locationTableNode == nil {
		return nil, htmlutils.WrapError(ErrLocationTableNotFound, rootNode)
	}

	tBodyNode := htmlutils.BFS(locationTableNode, matchTbody)
	if tBodyNode == nil {
		return nil, htmlutils.WrapError(ErrTbodyNotFound, rootNode)
	}

	rowNodes := htmlutils.DirectChildren(tBodyNode, matchTr)
	hostToData = make(map[string]serverData, len(rowNodes))

	for _, rowNode := range rowNodes {
		hostname, data, err := parseRowNode(rowNode)
		if err != nil {
			return nil, fmt.Errorf("parsing row node: %w", err)
		}
		hostToData[hostname] = data
	}

	return hostToData, nil
}

func parseRowNode(rowNode *html.Node) (hostname string, data serverData, err error) {
	columnIndex := 0
	const (
		columnIndexContinent = 0
		columnIndexCountry   = 1
		columnIndexCity      = 2
		columnIndexConfig    = 3
	)
	for cellNode := rowNode.FirstChild; cellNode != nil; cellNode = cellNode.NextSibling {
		if cellNode.FirstChild == nil {
			continue
		}

		switch columnIndex {
		case columnIndexContinent:
			data.region = cellNode.FirstChild.Data
		case columnIndexCountry:
			data.country = cellNode.FirstChild.Data
		case columnIndexCity:
			data.city = cellNode.FirstChild.Data
		case columnIndexConfig:
			linkNodes := htmlutils.DirectChildren(cellNode, matchA)
			for _, linkNode := range linkNodes {
				if linkNode.FirstChild.Data != "OpenVPN" {
					continue
				}

				data.ovpnURL = htmlutils.Attribute(linkNode, "href")
				if data.ovpnURL == "" {
					return "", data, htmlutils.WrapError(ErrExtractOpenVPNURL, linkNode)
				}

				hostname, err = extractHostnameFromURL(data.ovpnURL)
				if err != nil {
					return "", data, fmt.Errorf("extracting hostname from url: %w", err)
				}

				break
			}
		}

		columnIndex++
		if columnIndex == columnIndexConfig+1 {
			break
		}
	}

	return hostname, data, nil
}

func matchLocationTable(rootNode *html.Node) (match bool) {
	return htmlutils.MatchID("location-table")(rootNode)
}

func matchTbody(locationTableNode *html.Node) (match bool) {
	return htmlutils.MatchData("tbody")(locationTableNode)
}

func matchTr(tbodyNode *html.Node) (match bool) {
	return htmlutils.MatchData("tr")(tbodyNode)
}

func matchA(cellNode *html.Node) (match bool) {
	return htmlutils.MatchData("a")(cellNode)
}

var serverNameRegex = regexp.MustCompile(`^.+\/(?P<serverName>.+)\.ovpn$`)

var (
	ErrExtractHostnameFromURL = errors.New("cannot extract hostname from url")
)

func extractHostnameFromURL(url string) (hostname string, err error) {
	matches := serverNameRegex.FindStringSubmatch(url)
	const minMatches = 2
	if len(matches) < minMatches {
		return "", fmt.Errorf("%w: %s has less than 2 matches for %s",
			ErrExtractHostnameFromURL, url, serverNameRegex)
	}
	hostname = matches[1]
	return hostname, nil
}
