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
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	ErrFormatNotRecognized       = errors.New("format is not recognized")
	ErrProviderUnspecified       = errors.New("VPN provider to format was not specified")
	ErrMultipleProvidersToFormat = errors.New("more than one VPN provider to format were specified")
)

func addProviderFlag(flagSet *flag.FlagSet, providerToFormat map[string]*bool,
	provider string, titleCaser cases.Caser) {
	boolPtr, ok := providerToFormat[provider]
	if !ok {
		panic(fmt.Sprintf("unknown provider in format map: %s", provider))
	}
	flagSet.BoolVar(boolPtr, provider, false, "Format "+titleCaser.String(provider)+" servers")
}

func (c *CLI) FormatServers(args []string) error {
	var format, output string
	allProviders := providers.All()
	providersToFormat := make(map[string]*bool, len(allProviders))
	for _, provider := range allProviders {
		providersToFormat[provider] = new(bool)
	}
	flagSet := flag.NewFlagSet("markdown", flag.ExitOnError)
	flagSet.StringVar(&format, "format", "markdown", "Format to use which can be: 'markdown'")
	flagSet.StringVar(&output, "output", "/dev/stdout", "Output file to write the formatted data to")
	titleCaser := cases.Title(language.English)
	for _, provider := range allProviders {
		addProviderFlag(flagSet, providersToFormat, provider, titleCaser)
	}
	if err := flagSet.Parse(args); err != nil {
		return err
	}

	if format != "markdown" {
		return fmt.Errorf("%w: %s", ErrFormatNotRecognized, format)
	}

	// Verify only one provider is set to be formatted.
	var providers []string
	for provider, formatPtr := range providersToFormat {
		if *formatPtr {
			providers = append(providers, provider)
		}
	}
	switch len(providers) {
	case 0:
		return ErrProviderUnspecified
	case 1:
	default:
		return fmt.Errorf("%w: %d specified: %s",
			ErrMultipleProvidersToFormat, len(providers),
			strings.Join(providers, ", "))
	}
	providerToFormat := providers[0]

	logger := newNoopLogger()
	storage, err := storage.New(logger, constants.ServersData)
	if err != nil {
		return fmt.Errorf("cannot create servers storage: %w", err)
	}

	formatted := storage.FormatToMarkdown(providerToFormat)

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
