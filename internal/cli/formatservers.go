package cli

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/constants/providers"
	"github.com/qdm12/gluetun/internal/storage"
)

type ServersFormatter interface {
	FormatServers(args []string) error
}

var (
	ErrFormatNotRecognized = errors.New("format is not recognized")
	ErrProviderUnspecified = errors.New("VPN provider to format was not specified")
)

func addProviderFlag(flagSet *flag.FlagSet,
	providerToFormat map[string]*bool, provider string) {
	boolPtr, ok := providerToFormat[provider]
	if !ok {
		panic(fmt.Sprintf("unknown provider in format map: %s", provider))
	}
	flagSet.BoolVar(boolPtr, provider, false, "Format "+strings.Title(provider)+" servers")
}

func getFormatForProvider(providerToFormat map[string]*bool, provider string) (format bool) {
	formatPtr, ok := providerToFormat[provider]
	if !ok {
		panic(fmt.Sprintf("unknown provider in format map: %s", provider))
	}
	return *formatPtr
}

func (c *CLI) FormatServers(args []string) error {
	var format, output string
	allProviders := providers.All()
	providersToFormat := make(map[string]*bool, len(allProviders))
	for _, provider := range allProviders {
		value := false
		providersToFormat[provider] = &value
	}
	flagSet := flag.NewFlagSet("markdown", flag.ExitOnError)
	flagSet.StringVar(&format, "format", "markdown", "Format to use which can be: 'markdown'")
	flagSet.StringVar(&output, "output", "/dev/stdout", "Output file to write the formatted data to")
	for _, provider := range allProviders {
		addProviderFlag(flagSet, providersToFormat, provider)
	}
	if err := flagSet.Parse(args); err != nil {
		return err
	}

	if format != "markdown" {
		return fmt.Errorf("%w: %s", ErrFormatNotRecognized, format)
	}

	logger := newNoopLogger()
	storage, err := storage.New(logger, constants.ServersData)
	if err != nil {
		return fmt.Errorf("cannot create servers storage: %w", err)
	}
	currentServers := storage.GetServers()

	var formatted string
	switch {
	case getFormatForProvider(providersToFormat, providers.Cyberghost):
		formatted = currentServers.Cyberghost.ToMarkdown(providers.Cyberghost)
	case getFormatForProvider(providersToFormat, providers.Expressvpn):
		formatted = currentServers.Expressvpn.ToMarkdown(providers.Expressvpn)
	case getFormatForProvider(providersToFormat, providers.Fastestvpn):
		formatted = currentServers.Fastestvpn.ToMarkdown(providers.Fastestvpn)
	case getFormatForProvider(providersToFormat, providers.HideMyAss):
		formatted = currentServers.HideMyAss.ToMarkdown(providers.HideMyAss)
	case getFormatForProvider(providersToFormat, providers.Ipvanish):
		formatted = currentServers.Ipvanish.ToMarkdown(providers.Ipvanish)
	case getFormatForProvider(providersToFormat, providers.Ivpn):
		formatted = currentServers.Ivpn.ToMarkdown(providers.Ivpn)
	case getFormatForProvider(providersToFormat, providers.Mullvad):
		formatted = currentServers.Mullvad.ToMarkdown(providers.Mullvad)
	case getFormatForProvider(providersToFormat, providers.Nordvpn):
		formatted = currentServers.Nordvpn.ToMarkdown(providers.Nordvpn)
	case getFormatForProvider(providersToFormat, providers.Perfectprivacy):
		formatted = currentServers.Perfectprivacy.ToMarkdown(providers.Perfectprivacy)
	case getFormatForProvider(providersToFormat, providers.PrivateInternetAccess):
		formatted = currentServers.Pia.ToMarkdown(providers.PrivateInternetAccess)
	case getFormatForProvider(providersToFormat, providers.Privado):
		formatted = currentServers.Privado.ToMarkdown(providers.Privado)
	case getFormatForProvider(providersToFormat, providers.Privatevpn):
		formatted = currentServers.Privatevpn.ToMarkdown(providers.Privatevpn)
	case getFormatForProvider(providersToFormat, providers.Protonvpn):
		formatted = currentServers.Protonvpn.ToMarkdown(providers.Protonvpn)
	case getFormatForProvider(providersToFormat, providers.Purevpn):
		formatted = currentServers.Purevpn.ToMarkdown(providers.Purevpn)
	case getFormatForProvider(providersToFormat, providers.Surfshark):
		formatted = currentServers.Surfshark.ToMarkdown(providers.Surfshark)
	case getFormatForProvider(providersToFormat, providers.Torguard):
		formatted = currentServers.Torguard.ToMarkdown(providers.Torguard)
	case getFormatForProvider(providersToFormat, providers.VPNUnlimited):
		formatted = currentServers.VPNUnlimited.ToMarkdown(providers.VPNUnlimited)
	case getFormatForProvider(providersToFormat, providers.Vyprvpn):
		formatted = currentServers.Vyprvpn.ToMarkdown(providers.Vyprvpn)
	case getFormatForProvider(providersToFormat, providers.Wevpn):
		formatted = currentServers.Wevpn.ToMarkdown(providers.Wevpn)
	case getFormatForProvider(providersToFormat, providers.Windscribe):
		formatted = currentServers.Windscribe.ToMarkdown(providers.Windscribe)
	default:
		return ErrProviderUnspecified
	}

	output = filepath.Clean(output)
	file, err := os.OpenFile(output, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("cannot open output file: %w", err)
	}

	_, err = fmt.Fprint(file, formatted)
	if err != nil {
		_ = file.Close()
		return fmt.Errorf("cannot write to output file: %w", err)
	}

	err = file.Close()
	if err != nil {
		return fmt.Errorf("cannot close output file: %w", err)
	}

	return nil
}
