package settings

import (
	"fmt"
	"strings"
	"time"

	"github.com/qdm12/gluetun/internal/params"
)

type Updater struct {
	Period     time.Duration `json:"period"`
	DNSAddress string        `json:"dns_address"`
	Cyberghost bool          `json:"cyberghost"`
	Mullvad    bool          `json:"mullvad"`
	Nordvpn    bool          `json:"nordvpn"`
	PIA        bool          `json:"pia"`
	Privado    bool          `json:"privado"`
	Purevpn    bool          `json:"purevpn"`
	Surfshark  bool          `json:"surfshark"`
	Vyprvpn    bool          `json:"vyprvpn"`
	Windscribe bool          `json:"windscribe"`
	// The two below should be used in CLI mode only
	Stdout bool `json:"-"` // in order to update constants file (maintainer side)
	CLI    bool `json:"-"`
}

// GetUpdaterSettings obtains the server updater settings using the params functions.
func GetUpdaterSettings(paramsReader params.Reader) (settings Updater, err error) {
	settings = Updater{
		Cyberghost: true,
		Mullvad:    true,
		Nordvpn:    true,
		PIA:        true,
		Purevpn:    true,
		Surfshark:  true,
		Vyprvpn:    true,
		Windscribe: true,
		Stdout:     false,
		CLI:        false,
		DNSAddress: "127.0.0.1",
	}
	settings.Period, err = paramsReader.GetUpdaterPeriod()
	if err != nil {
		return settings, err
	}
	return settings, nil
}

func (s *Updater) String() string {
	if s.Period == 0 {
		return "Server updater settings: disabled"
	}
	settingsList := []string{
		"Server updater settings:",
		fmt.Sprintf("Period: %s", s.Period),
	}
	return strings.Join(settingsList, "\n|--")
}
