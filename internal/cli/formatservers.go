package cli

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

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

func (c *CLI) FormatServers(args []string) error {
	var format, output string
	var cyberghost, expressvpn, fastestvpn, hideMyAss, ipvanish, ivpn, mullvad,
		nordvpn, perfectPrivacy, pia, privado, privatevpn, protonvpn, purevpn, surfshark,
		torguard, vpnUnlimited, vyprvpn, wevpn, windscribe bool
	flagSet := flag.NewFlagSet("markdown", flag.ExitOnError)
	flagSet.StringVar(&format, "format", "markdown", "Format to use which can be: 'markdown'")
	flagSet.StringVar(&output, "output", "/dev/stdout", "Output file to write the formatted data to")
	flagSet.BoolVar(&cyberghost, providers.Cyberghost, false, "Format Cyberghost servers")
	flagSet.BoolVar(&expressvpn, providers.Expressvpn, false, "Format ExpressVPN servers")
	flagSet.BoolVar(&fastestvpn, providers.Fastestvpn, false, "Format FastestVPN servers")
	flagSet.BoolVar(&hideMyAss, providers.HideMyAss, false, "Format HideMyAss servers")
	flagSet.BoolVar(&ipvanish, providers.Ipvanish, false, "Format IpVanish servers")
	flagSet.BoolVar(&ivpn, providers.Ivpn, false, "Format IVPN servers")
	flagSet.BoolVar(&mullvad, providers.Mullvad, false, "Format Mullvad servers")
	flagSet.BoolVar(&nordvpn, providers.Nordvpn, false, "Format Nordvpn servers")
	flagSet.BoolVar(&perfectPrivacy, providers.Perfectprivacy, false, "Format Perfect Privacy servers")
	flagSet.BoolVar(&pia, providers.PrivateInternetAccess, false, "Format Private Internet Access servers")
	flagSet.BoolVar(&privado, providers.Privado, false, "Format Privado servers")
	flagSet.BoolVar(&privatevpn, providers.Privatevpn, false, "Format Private VPN servers")
	flagSet.BoolVar(&protonvpn, providers.Protonvpn, false, "Format Protonvpn servers")
	flagSet.BoolVar(&purevpn, providers.Purevpn, false, "Format Purevpn servers")
	flagSet.BoolVar(&surfshark, providers.Surfshark, false, "Format Surfshark servers")
	flagSet.BoolVar(&torguard, providers.Torguard, false, "Format Torguard servers")
	flagSet.BoolVar(&vpnUnlimited, providers.VPNUnlimited, false, "Format VPN Unlimited servers")
	flagSet.BoolVar(&vyprvpn, providers.Vyprvpn, false, "Format Vyprvpn servers")
	flagSet.BoolVar(&wevpn, providers.Wevpn, false, "Format WeVPN servers")
	flagSet.BoolVar(&windscribe, providers.Windscribe, false, "Format Windscribe servers")
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
		formatted = currentServers.Cyberghost.ToMarkdown(providers.Cyberghost)
	case expressvpn:
		formatted = currentServers.Expressvpn.ToMarkdown(providers.Expressvpn)
	case fastestvpn:
		formatted = currentServers.Fastestvpn.ToMarkdown(providers.Fastestvpn)
	case hideMyAss:
		formatted = currentServers.HideMyAss.ToMarkdown(providers.HideMyAss)
	case ipvanish:
		formatted = currentServers.Ipvanish.ToMarkdown(providers.Ipvanish)
	case ivpn:
		formatted = currentServers.Ivpn.ToMarkdown(providers.Ivpn)
	case mullvad:
		formatted = currentServers.Mullvad.ToMarkdown(providers.Mullvad)
	case nordvpn:
		formatted = currentServers.Nordvpn.ToMarkdown(providers.Nordvpn)
	case perfectPrivacy:
		formatted = currentServers.Perfectprivacy.ToMarkdown(providers.Perfectprivacy)
	case pia:
		formatted = currentServers.Pia.ToMarkdown(providers.PrivateInternetAccess)
	case privado:
		formatted = currentServers.Privado.ToMarkdown(providers.Privado)
	case privatevpn:
		formatted = currentServers.Privatevpn.ToMarkdown(providers.Privatevpn)
	case protonvpn:
		formatted = currentServers.Protonvpn.ToMarkdown(providers.Protonvpn)
	case purevpn:
		formatted = currentServers.Purevpn.ToMarkdown(providers.Purevpn)
	case surfshark:
		formatted = currentServers.Surfshark.ToMarkdown(providers.Surfshark)
	case torguard:
		formatted = currentServers.Torguard.ToMarkdown(providers.Torguard)
	case vpnUnlimited:
		formatted = currentServers.VPNUnlimited.ToMarkdown(providers.VPNUnlimited)
	case vyprvpn:
		formatted = currentServers.Vyprvpn.ToMarkdown(providers.Vyprvpn)
	case wevpn:
		formatted = currentServers.Wevpn.ToMarkdown(providers.Wevpn)
	case windscribe:
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
