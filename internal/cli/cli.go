package cli

import (
	"flag"
	"fmt"
	"strings"

	"net"

	"github.com/qdm12/golibs/files"
)

func ClientKey(args []string) error {
	flagSet := flag.NewFlagSet("clientkey", flag.ExitOnError)
	filepath := flagSet.String("path", "/client.key", "file path to the client.key file")
	if err := flagSet.Parse(args); err != nil {
		return err
	}
	fileManager := files.NewFileManager()
	data, err := fileManager.ReadFile(*filepath)
	if err != nil {
		return err
	}
	s := string(data)
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.TrimPrefix(s, "-----BEGIN PRIVATE KEY-----")
	s = strings.TrimSuffix(s, "-----END PRIVATE KEY-----")
	fmt.Println(s)
	return nil
}

func HealthCheck() error {
	ips, err := net.LookupIP("github.com")
	if err != nil {
		return fmt.Errorf("cannot resolve github.com (%s)", err)
	} else if len(ips) == 0 {
		return fmt.Errorf("resolved no IP addresses for github.com")
	}
	return nil
}
