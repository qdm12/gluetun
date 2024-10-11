package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

func (c *CLI) ClientKey(args []string) error {
	flagSet := flag.NewFlagSet("clientkey", flag.ExitOnError)
	const openVPNClientKeyPath = "/gluetun/client.key" // TODO deduplicate?
	filepath := flagSet.String("path", openVPNClientKeyPath, "file path to the client.key file")
	if err := flagSet.Parse(args); err != nil {
		return err
	}
	file, err := os.OpenFile(*filepath, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	data, err := io.ReadAll(file)
	if err != nil {
		_ = file.Close()
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	s := string(data)
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "\r", "")
	s = strings.TrimPrefix(s, "-----BEGIN PRIVATE KEY-----")
	s = strings.TrimSuffix(s, "-----END PRIVATE KEY-----")
	fmt.Println(s)
	return nil
}
