package updater

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/qdm12/gluetun/internal/models"
)

func fetchServers(ctx context.Context, client *http.Client) (
	servers []models.Server, err error) {
	data, err := fetchHTML(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("cannot fetch HTML code: %w", err)
	}

	servers = parseHTML(string(data))
	return servers, nil
}

var ErrHTTPStatusCode = errors.New("HTTP status code is not OK")

func fetchHTML(ctx context.Context, client *http.Client) (data []byte, err error) {
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

	data, err = io.ReadAll(response.Body)
	if err != nil {
		_ = response.Body.Close()
		return nil, err
	}

	err = response.Body.Close()
	if err != nil {
		return nil, err
	}

	return data, nil
}

func parseHTML(html string) (servers []models.Server) {
	// Remove consecutive empty lines
	for strings.Contains(html, "\n\n") {
		html = strings.ReplaceAll(html, "\n\n", "\n")
	}

	var block string
	var readingBlock bool
	for _, line := range strings.Split(html, "\n") {
		line = strings.TrimSpace(line)

		const (
			blockStartString = `<dl class="grid__i">`
			blockEndString   = `</dl>`
		)
		switch {
		case strings.HasPrefix(line, blockStartString):
			readingBlock = true
		case strings.HasPrefix(line, blockEndString):
			readingBlock = false
			server, ok := extractFromHTMLBlock(block)
			if ok {
				servers = append(servers, server)
			}
			block = ""
		case readingBlock &&
			!strings.Contains(line, "</svg>"): // ignore SVG element lines
			block += line + "\n"
		}
	}

	return servers
}

var (
	hostRegex    = regexp.MustCompile(`<dt>[a-z]+[0-9]*<span`)
	cityRegex    = regexp.MustCompile(`<div><span>City:</span> <strong>.+?</strong></div>`)
	regionRegex  = regexp.MustCompile(`<div><span>Region:</span> <strong>.+?</strong></div>`)
	premiumRegex = regexp.MustCompile(`<div><span>Premium:</span> <strong>[a-zA-Z]+?</strong></div>`)
)

func extractFromHTMLBlock(htmlBlock string) (server models.Server, ok bool) {
	htmlBlock = strings.ReplaceAll(htmlBlock, "\n", "")

	host := regexTrimPrefixSuffix(htmlBlock, hostRegex,
		"<dt>", "<span")
	if host == "" {
		return server, false
	}
	server.Hostname = host + ".isponeder.com"

	server.City = regexTrimPrefixSuffix(htmlBlock, cityRegex,
		"<div><span>City:</span> <strong>", "</strong></div>")

	server.Region = regexTrimPrefixSuffix(htmlBlock, regionRegex,
		"<div><span>Region:</span> <strong>", "</strong></div>")

	premiumString := regexTrimPrefixSuffix(htmlBlock, premiumRegex,
		"<div><span>Premium:</span> <strong>",
		"</strong></div>")
	server.Premium = strings.EqualFold(premiumString, "yes")

	return server, true
}

func regexTrimPrefixSuffix(s string, regex *regexp.Regexp,
	prefix, suffix string) (result string) {
	result = regex.FindString(s)
	result = strings.TrimPrefix(result, prefix)
	result = strings.TrimSuffix(result, suffix)
	return result
}
