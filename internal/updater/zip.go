package updater

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
)

func fetchAndExtractFiles(ctx context.Context, client *http.Client, urls ...string) (
	contents map[string][]byte, err error) {
	contents = make(map[string][]byte)
	for _, url := range urls {
		request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}

		response, err := client.Do(request)
		if err != nil {
			return nil, err
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("%w: %s for %s", ErrHTTPStatusCodeNotOK, response.Status, url)
		}

		b, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}

		if err := response.Body.Close(); err != nil {
			return nil, err
		}

		newContents, err := zipExtractAll(b)
		if err != nil {
			return nil, err
		}
		for fileName, content := range newContents {
			contents[fileName] = content
		}
	}
	return contents, nil
}

func zipExtractAll(zipBytes []byte) (contents map[string][]byte, err error) {
	r, err := zip.NewReader(bytes.NewReader(zipBytes), int64(len(zipBytes)))
	if err != nil {
		return nil, err
	}
	contents = map[string][]byte{}
	for _, zf := range r.File {
		fileName := filepath.Base(zf.Name)
		if !strings.HasSuffix(fileName, ".ovpn") {
			continue
		}
		f, err := zf.Open()
		if err != nil {
			return nil, err
		}
		defer f.Close()
		contents[fileName], err = ioutil.ReadAll(f)
		if err != nil {
			return nil, err
		}
		if err := f.Close(); err != nil {
			return nil, err
		}
	}
	return contents, nil
}
