package updater

import (
	"context"
	"io"
	"net/http"
	"regexp"
	"strings"
)

var indexOpenvpnLinksRegex = regexp.MustCompile(`<a[ ]+href=".+\.ovpn">.+\.ovpn</a>`)

func fetchIndex(ctx context.Context, client *http.Client, indexURL string) (
	openvpnURLs []string, err error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, indexURL, nil)
	if err != nil {
		return nil, err
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	htmlCode, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if !strings.HasSuffix(indexURL, "/") {
		indexURL += "/"
	}

	lines := strings.Split(string(htmlCode), "\n")
	for _, line := range lines {
		found := indexOpenvpnLinksRegex.FindString(line)
		if len(found) == 0 {
			continue
		}
		const prefix = `.ovpn">`
		const suffix = `</a>`
		startIndex := strings.Index(found, prefix) + len(prefix)
		endIndex := strings.Index(found, suffix)
		filename := found[startIndex:endIndex]
		openvpnURL := indexURL + filename
		if !strings.HasSuffix(openvpnURL, ".ovpn") {
			continue
		}
		openvpnURLs = append(openvpnURLs, openvpnURL)
	}

	return openvpnURLs, nil
}
