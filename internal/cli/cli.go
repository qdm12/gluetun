package cli

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/params"
	"github.com/qdm12/gluetun/internal/provider"
	"github.com/qdm12/gluetun/internal/settings"
	"github.com/qdm12/gluetun/internal/storage"
	"github.com/qdm12/gluetun/internal/updater"
	"github.com/qdm12/golibs/files"
	"github.com/qdm12/golibs/logging"
	"golang.org/x/net/context/ctxhttp"
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

func HealthCheck(ctx context.Context, client *http.Client) error {
	const url = "http://localhost:8000/health"
	response, err := ctxhttp.Get(ctx, client, url)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusOK {
		return nil
	}
	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	return fmt.Errorf("HTTP status code %s with message: %s", response.Status, string(b))
}

func OpenvpnConfig() error {
	logger, err := logging.NewLogger(logging.ConsoleEncoding, logging.InfoLevel)
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
	providerConf := provider.New(allSettings.OpenVPN.Provider.Name, allServers, time.Now)
	connection, err := providerConf.GetOpenVPNConnection(allSettings.OpenVPN.Provider.ServerSelection)
	if err != nil {
		return err
	}
	lines := providerConf.BuildConf(
		connection,
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
	options := updater.Options{CLI: true}
	var flushToFile bool
	flagSet := flag.NewFlagSet("update", flag.ExitOnError)
	flagSet.BoolVar(&flushToFile, "file", false, "Write results to /gluetun/servers.json (for end users)")
	flagSet.BoolVar(&options.Stdout, "stdout", false, "Write results to console to modify the program (for maintainers)")
	flagSet.StringVar(&options.DNSAddress, "dns", "1.1.1.1", "DNS resolver address to use")
	flagSet.BoolVar(&options.Cyberghost, "cyberghost", false, "Update Cyberghost servers")
	flagSet.BoolVar(&options.Mullvad, "mullvad", false, "Update Mullvad servers")
	flagSet.BoolVar(&options.Nordvpn, "nordvpn", false, "Update Nordvpn servers")
	flagSet.BoolVar(&options.PIA, "pia", false, "Update Private Internet Access post-summer 2020 servers")
	flagSet.BoolVar(&options.PIAold, "piaold", false, "Update Private Internet Access pre-summer 2020 servers")
	flagSet.BoolVar(&options.Purevpn, "purevpn", false, "Update Purevpn servers")
	flagSet.BoolVar(&options.Surfshark, "surfshark", false, "Update Surfshark servers")
	flagSet.BoolVar(&options.Vyprvpn, "vyprvpn", false, "Update Vyprvpn servers")
	flagSet.BoolVar(&options.Windscribe, "windscribe", false, "Update Windscribe servers")
	if err := flagSet.Parse(args); err != nil {
		return err
	}
	logger, err := logging.NewLogger(logging.ConsoleEncoding, logging.InfoLevel)
	if err != nil {
		return err
	}
	if !flushToFile && !options.Stdout {
		return fmt.Errorf("at least one of -file or -stdout must be specified")
	}
	ctx := context.Background()
	httpClient := &http.Client{Timeout: 10 * time.Second}
	storage := storage.New(logger)
	const writeSync = false
	currentServers, err := storage.SyncServers(constants.GetAllServers(), writeSync)
	if err != nil {
		return fmt.Errorf("cannot update servers: %w", err)
	}
	updater := updater.New(options, httpClient, currentServers, logger)
	allServers, err := updater.UpdateServers(ctx)
	if err != nil {
		return err
	}
	if flushToFile {
		if err := storage.FlushToFile(allServers); err != nil {
			return fmt.Errorf("cannot update servers: %w", err)
		}
	}

	return nil
}
