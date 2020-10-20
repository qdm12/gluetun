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

	"github.com/qdm12/golibs/network"
)

func fetchAndExtractFiles(ctx context.Context, client network.Client, urls ...string) (
	contents map[string][]byte, err error) {
	contents = make(map[string][]byte)
	for _, url := range urls {
		zipBytes, status, err := client.Get(ctx, url)
		if err != nil {
			return nil, err
		} else if status != http.StatusOK {
			return nil, fmt.Errorf("Getting %s results in HTTP status code %d", url, status)
		}
		newContents, err := zipExtractAll(zipBytes)
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
