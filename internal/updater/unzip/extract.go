package unzip

import (
	"archive/zip"
	"bytes"
	"io/ioutil"
	"path/filepath"
	"strings"
)

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
