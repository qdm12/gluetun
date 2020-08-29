package cli

import (
	"flag"
	"fmt"
	"net/http"
	"strings"
	"time"

	"net"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/params"
	"github.com/qdm12/gluetun/internal/provider"
	"github.com/qdm12/gluetun/internal/settings"
	"github.com/qdm12/gluetun/internal/storage"
	"github.com/qdm12/gluetun/internal/updater"
	"github.com/qdm12/golibs/files"
	"github.com/qdm12/golibs/logging"
)

func ClientKey(args []string) error {
	flagSet := flag.NewFlagSet("clientkey", flag.ExitOnError)
	filepath := flagSet.String("path", "/files/client.key", "file path to the client.key file")
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
	s = strings.ReplaceAll(s, "\r", "")
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

func OpenvpnConfig() error {
	logger, err := logging.NewLogger(logging.ConsoleEncoding, logging.InfoLevel, -1)
	if err != nil {
		return err
	}
	paramsReader := params.NewReader(logger, files.NewFileManager())
	allSettings, err := settings.GetAllSettings(paramsReader)
	if err != nil {
		return err
	}
	allServers, err := storage.New(logger).SyncServers(constants.GetAllServers(), false)
	if err != nil {
		return err
	}
	providerConf := provider.New(allSettings.OpenVPN.Provider.Name, allServers)
	connections, err := providerConf.GetOpenVPNConnections(allSettings.OpenVPN.Provider.ServerSelection)
	if err != nil {
		return err
	}
	lines := providerConf.BuildConf(
		connections,
		allSettings.OpenVPN.Verbosity,
		allSettings.System.UID,
		allSettings.System.GID,
		allSettings.OpenVPN.Root,
		allSettings.OpenVPN.Cipher,
		allSettings.OpenVPN.Auth,
		allSettings.OpenVPN.Provider.ExtraConfigOptions,
	)
	fmt.Println(strings.Join(lines, "\n"))
	return nil
}

func Update(args []string) error {
	var options updater.Options
	flagSet := flag.NewFlagSet("update", flag.ExitOnError)
	flagSet.BoolVar(&options.File, "file", false, "Write results to /gluetun/servers.json (for end users)")
	flagSet.BoolVar(&options.Stdout, "stdout", false, "Write results to console to modify the program (for maintainers)")
	flagSet.BoolVar(&options.PIA, "pia", false, "Update Private Internet Access post-summer 2020 servers")
	flagSet.BoolVar(&options.PIAold, "piaold", false, "Update Private Internet Access pre-summer 2020 servers")
	flagSet.BoolVar(&options.Mullvad, "mullvad", false, "Update Mullvad servers")
	if err := flagSet.Parse(args); err != nil {
		return err
	}
	logger, err := logging.NewLogger(logging.ConsoleEncoding, logging.InfoLevel, -1)
	if err != nil {
		return err
	}
	if !options.File && !options.Stdout {
		return fmt.Errorf("at least one of -file or -stdout must be specified")
	}
	httpClient := &http.Client{Timeout: 10 * time.Second}
	storage := storage.New(logger)
	updater := updater.New(storage, httpClient)
	if err := updater.UpdateServers(options); err != nil {
		return err
	}
	return nil
}
