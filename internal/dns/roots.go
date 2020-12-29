package dns

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/qdm12/gluetun/internal/constants"
)

func (c *configurator) DownloadRootHints(ctx context.Context, puid, pgid int) error {
	return c.downloadAndSave(ctx, "root hints",
		string(constants.NamedRootURL), string(constants.RootHints), puid, pgid)
}

func (c *configurator) DownloadRootKey(ctx context.Context, puid, pgid int) error {
	return c.downloadAndSave(ctx, "root key",
		string(constants.RootKeyURL), string(constants.RootKey), puid, pgid)
}

func (c *configurator) downloadAndSave(ctx context.Context, logName, url, filepath string, puid, pgid int) error {
	c.logger.Info("downloading %s from %s", logName, url)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	response, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("%w from %s: %s", ErrBadStatusCode, url, response.Status)
	}

	file, err := c.openFile(filepath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0400)
	if err != nil {
		return err
	}

	_, err = io.Copy(file, response.Body)
	if err != nil {
		_ = file.Close()
		return err
	}

	err = file.Chown(puid, pgid)
	if err != nil {
		_ = file.Close()
		return err
	}

	return file.Close()
}
