package cli

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/storage"
)

type ServersFormatter interface {
	FormatServers(args []string) error
}

var (
	ErrFormatNotRecognized = errors.New("format is not recognized")
	ErrProviderUnspecified = errors.New("VPN provider to format was not specified")
)

func (c *CLI) FormatServers(args []string) error {
	var format, output string
	var cyberghost, expressvpn, fastestvpn, hideMyAss, ipvanish, ivpn, mullvad,
		nordvpn, perfectPrivacy, pia, privado, privatevpn, protonvpn, purevpn, surfshark,
		torguard, vpnUnlimited, vyprvpn, wevpn, windscribe bool
	flagSet := flag.NewFlagSet("markdown", flag.ExitOnError)
	flagSet.StringVar(&format, "format", "markdown", "Format to use which can be: 'markdown'")
	flagSet.StringVar(&output, "output", "/dev/stdout", "Output file to write the formatted data to")
	flagSet.BoolVar(&cyberghost, "cyberghost", false, "Format Cyberghost servers")
	flagSet.BoolVar(&expressvpn, "expressvpn", false, "Format ExpressVPN servers")
	flagSet.BoolVar(&fastestvpn, "fastestvpn", false, "Format FastestVPN servers")
	flagSet.BoolVar(&hideMyAss, "hidemyass", false, "Format HideMyAss servers")
	flagSet.BoolVar(&ipvanish, "ipvanish", false, "Format IpVanish servers")
	flagSet.BoolVar(&ivpn, "ivpn", false, "Format IVPN servers")
	flagSet.BoolVar(&mullvad, "mullvad", false, "Format Mullvad servers")
	flagSet.BoolVar(&nordvpn, "nordvpn", false, "Format Nordvpn servers")
	flagSet.BoolVar(&perfectPrivacy, "perfectprivacy", false, "Format Perfect Privacy servers")
	flagSet.BoolVar(&pia, "pia", false, "Format Private Internet Access servers")
	flagSet.BoolVar(&privado, "privado", false, "Format Privado servers")
	flagSet.BoolVar(&privatevpn, "privatevpn", false, "Format Private VPN servers")
	flagSet.BoolVar(&protonvpn, "protonvpn", false, "Format Protonvpn servers")
	flagSet.BoolVar(&purevpn, "purevpn", false, "Format Purevpn servers")
	flagSet.BoolVar(&surfshark, "surfshark", false, "Format Surfshark servers")
	flagSet.BoolVar(&torguard, "torguard", false, "Format Torguard servers")
	flagSet.BoolVar(&vpnUnlimited, "vpnunlimited", false, "Format VPN Unlimited servers")
	flagSet.BoolVar(&vyprvpn, "vyprvpn", false, "Format Vyprvpn servers")
	flagSet.BoolVar(&wevpn, "wevpn", false, "Format WeVPN servers")
	flagSet.BoolVar(&windscribe, "windscribe", false, "Format Windscribe servers")
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
	case cyberghost:
		formatted = currentServers.Cyberghost.ToMarkdown()
	case expressvpn:
		formatted = currentServers.Expressvpn.ToMarkdown()
	case fastestvpn:
		formatted = currentServers.Fastestvpn.ToMarkdown()
	case hideMyAss:
		formatted = currentServers.HideMyAss.ToMarkdown()
	case ipvanish:
		formatted = currentServers.Ipvanish.ToMarkdown()
	case ivpn:
		formatted = currentServers.Ivpn.ToMarkdown()
	case mullvad:
		formatted = currentServers.Mullvad.ToMarkdown()
	case nordvpn:
		formatted = currentServers.Nordvpn.ToMarkdown()
	case perfectPrivacy:
		formatted = currentServers.Perfectprivacy.ToMarkdown()
	case pia:
		formatted = currentServers.Pia.ToMarkdown()
	case privado:
		formatted = currentServers.Privado.ToMarkdown()
	case privatevpn:
		formatted = currentServers.Privatevpn.ToMarkdown()
	case protonvpn:
		formatted = currentServers.Protonvpn.ToMarkdown()
	case purevpn:
		formatted = currentServers.Purevpn.ToMarkdown()
	case surfshark:
		formatted = currentServers.Surfshark.ToMarkdown()
	case torguard:
		formatted = currentServers.Torguard.ToMarkdown()
	case vpnUnlimited:
		formatted = currentServers.VPNUnlimited.ToMarkdown()
	case vyprvpn:
		formatted = currentServers.Vyprvpn.ToMarkdown()
	case wevpn:
		formatted = currentServers.Wevpn.ToMarkdown()
	case windscribe:
		formatted = currentServers.Windscribe.ToMarkdown()
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
