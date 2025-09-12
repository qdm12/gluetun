package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io/fs"
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
	ErrProviderUnspecified       = errors.New("VPN provider to format was not specified")
	ErrMultipleProvidersToFormat = errors.New("more than one VPN provider to format were specified")
)

func addProviderFlag(flagSet *flag.FlagSet, providerToFormat map[string]*bool,
	provider string, titleCaser cases.Caser,
) {
	boolPtr, ok := providerToFormat[provider]
	if !ok {
		panic(fmt.Sprintf("unknown provider in format map: %s", provider))
	}
	flagSet.BoolVar(boolPtr, provider, false, "Format "+titleCaser.String(provider)+" servers")
}

type FormatServersCommand struct {
	args []string
}

func NewFormatServersCommand(args []string) *FormatServersCommand {
	return &FormatServersCommand{
		args: args,
	}
}

func (c *FormatServersCommand) Name() string {
	return "format-servers"
}

func (c *FormatServersCommand) Description() string {
	return "Format the servers data into a Markdown table"
}

func (c *FormatServersCommand) Run(_ context.Context) error {
	var format, output string
	allProviders := providers.All()
	allProviderFlags := make([]string, len(allProviders))
	for i, provider := range allProviders {
		allProviderFlags[i] = strings.ReplaceAll(provider, " ", "-")
	}

	providersToFormat := make(map[string]*bool, len(allProviders))
	for _, provider := range allProviderFlags {
		providersToFormat[provider] = new(bool)
	}
	flagSet := flag.NewFlagSet("format-servers", flag.ExitOnError)
	flagSet.StringVar(&format, "format", "markdown", "Format to use which can be: 'markdown' or 'json'")
	flagSet.StringVar(&output, "output", "/dev/stdout", "Output file to write the formatted data to")
	titleCaser := cases.Title(language.English)
	for _, provider := range allProviderFlags {
		addProviderFlag(flagSet, providersToFormat, provider, titleCaser)
	}
	if err := flagSet.Parse(c.args[2:]); err != nil {
		return err
	}

	// Note the format is validated by storage.Format

	// Verify only one provider is set to be formatted.
	var providers []string
	for provider, formatPtr := range providersToFormat {
		if *formatPtr {
			providers = append(providers, provider)
		}
	}
	switch len(providers) {
	case 0:
		return fmt.Errorf("%w", ErrProviderUnspecified)
	case 1:
	default:
		return fmt.Errorf("%w: %d specified: %s",
			ErrMultipleProvidersToFormat, len(providers),
			strings.Join(providers, ", "))
	}

	var providerToFormat string
	for _, providerToFormat = range allProviders {
		if strings.ReplaceAll(providerToFormat, " ", "-") == providers[0] {
			break
		}
	}

	logger := newNoopLogger()
	storage, err := storage.New(logger, constants.ServersData)
	if err != nil {
		return fmt.Errorf("creating servers storage: %w", err)
	}

	formatted, err := storage.Format(providerToFormat, format)
	if err != nil {
		return fmt.Errorf("formatting servers: %w", err)
	}

	output = filepath.Clean(output)
	const permission = fs.FileMode(0o644)
	file, err := os.OpenFile(output, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, permission)
	if err != nil {
		return fmt.Errorf("opening output file: %w", err)
	}

	_, err = fmt.Fprint(file, formatted)
	if err != nil {
		_ = file.Close()
		return fmt.Errorf("writing to output file: %w", err)
	}

	err = file.Close()
	if err != nil {
		return fmt.Errorf("closing output file: %w", err)
	}

	return nil
}
