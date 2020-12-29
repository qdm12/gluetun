package dns

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/qdm12/gluetun/internal/constants"
)

func (c *configurator) DownloadRootHints(ctx context.Context, uid, gid int) error {
	return c.downloadAndSave(ctx, "root hints",
		string(constants.NamedRootURL), string(constants.RootHints), uid, gid)
}

func (c *configurator) DownloadRootKey(ctx context.Context, uid, gid int) error {
	return c.downloadAndSave(ctx, "root key",
		string(constants.RootKeyURL), string(constants.RootKey), uid, gid)
}

func (c *configurator) downloadAndSave(ctx context.Context, logName, url, filepath string, uid, gid int) error {
	c.logger.Info("downloading %s from %s", logName, url)
	content, status, err := c.client.Get(ctx, url)
	if err != nil {
		return err
	} else if status != http.StatusOK {
		return fmt.Errorf("HTTP status code is %d for %s", status, url)
	}

	file, err := c.openFile(filepath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0400)
	if err != nil {
		return err
	}

	_, err = file.Write(content)
	if err != nil {
		_ = file.Close()
		return err
	}

	err = file.Chown(uid, gid)
	if err != nil {
		_ = file.Close()
		return err
	}

	return file.Close()
}
