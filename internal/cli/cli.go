package cli

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/qdm12/golibs/files"
	"github.com/qdm12/golibs/network/connectivity"
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
	// DNS, HTTP and HTTPs check on github.com
	connectivity := connectivity.NewConnectivity(3 * time.Second)
	errs := connectivity.Checks("github.com")
	if len(errs) > 0 {
		var errsStr []string
		for _, err := range errs {
			errsStr = append(errsStr, err.Error())
		}
		return fmt.Errorf("Multiple errors: %s", strings.Join(errsStr, "; "))
	}
	// TODO check IP address is in the right region
	return nil
}
