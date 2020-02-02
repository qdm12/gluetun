package dns

import (
	"fmt"

	"github.com/qdm12/private-internet-access-docker/internal/constants"
)

func (c *configurator) DownloadRootHints() error {
	content, status, err := c.client.GetContent(string(constants.NamedRootURL))
	if err != nil {
		return err
	} else if status != 200 {
		return fmt.Errorf("HTTP status code is %d for %s", status, constants.NamedRootURL)
	}
	return c.fileManager.WriteToFile(string(constants.RootHints), content)
}

func (c *configurator) DownloadRootKey() error {
	content, status, err := c.client.GetContent(string(constants.NamedRootURL))
	if err != nil {
		return err
	} else if status != 200 {
		return fmt.Errorf("HTTP status code is %d for %s", status, constants.RootKeyURL)
	}
	return c.fileManager.WriteToFile(string(constants.RootKey), content)
}
