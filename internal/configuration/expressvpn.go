package configuration

import (
	"errors"
	"fmt"

	"github.com/qdm12/gluetun/internal/constants"
)

var errProtocolNotSupported = errors.New("protocol is not supported")

func (settings *Provider) readExpressvpn(r reader) (err error) {
	settings.Name = constants.Expressvpn
	servers := r.servers.GetExpressvpn()

	settings.ServerSelection.TargetIP, err = readTargetIP(r.env)
	if err != nil {
		return err
	}

	settings.ServerSelection.Hostnames, err = r.env.CSVInside("SERVER_HOSTNAME",
		constants.ExpressvpnHostnameChoices(servers))
	if err != nil {
		return fmt.Errorf("environment variable SERVER_HOSTNAME: %w", err)
	}

	settings.ServerSelection.Countries, err = r.env.CSVInside("COUNTRY", constants.ExpressvpnCountriesChoices(servers))
	if err != nil {
		return fmt.Errorf("environment variable COUNTRY: %w", err)
	}

	settings.ServerSelection.Cities, err = r.env.CSVInside("CITY", constants.ExpressvpnCitiesChoices(servers))
	if err != nil {
		return fmt.Errorf("environment variable CITY: %w", err)
	}

	tcp, _ := readProtocol(r.env)
	if tcp {
		return fmt.Errorf("%w: for provider %s", errProtocolNotSupported, constants.Expressvpn)
	}

	return nil
}
