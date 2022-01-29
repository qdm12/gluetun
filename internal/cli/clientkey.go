package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/qdm12/gluetun/internal/configuration/sources/files"
)

type ClientKeyFormatter interface {
	ClientKey(args []string) error
}

func (c *CLI) ClientKey(args []string) error {
	flagSet := flag.NewFlagSet("clientkey", flag.ExitOnError)
	filepath := flagSet.String("path", files.OpenVPNClientKeyPath, "file path to the client.key file")
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
	if err != nil {
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
