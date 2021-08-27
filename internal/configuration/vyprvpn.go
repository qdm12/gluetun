package configuration

import (
	"fmt"

	"github.com/qdm12/gluetun/internal/constants"
)

func (settings *Provider) readVyprvpn(r reader) (err error) {
	settings.Name = constants.Vyprvpn
	servers := r.servers.GetVyprvpn()

	settings.ServerSelection.TargetIP, err = readTargetIP(r.env)
	if err != nil {
		return err
	}

	settings.ServerSelection.Regions, err = r.env.CSVInside("REGION", constants.VyprvpnRegionChoices(servers))
	if err != nil {
		return fmt.Errorf("environment variable REGION: %w", err)
	}

	return nil
}
