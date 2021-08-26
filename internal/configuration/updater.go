package configuration

import (
	"fmt"
	"strings"
	"time"

	"github.com/qdm12/golibs/params"
)

type Updater struct {
	Period       time.Duration `json:"period"`
	DNSAddress   string        `json:"dns_address"`
	Cyberghost   bool          `json:"cyberghost"`
	Fastestvpn   bool          `json:"fastestvpn"`
	HideMyAss    bool          `json:"hidemyass"`
	Ipvanish     bool          `json:"ipvanish"`
	Ivpn         bool          `json:"ivpn"`
	Mullvad      bool          `json:"mullvad"`
	Nordvpn      bool          `json:"nordvpn"`
	PIA          bool          `json:"pia"`
	Privado      bool          `json:"privado"`
	Privatevpn   bool          `json:"privatevpn"`
	Protonvpn    bool          `json:"protonvpn"`
	Purevpn      bool          `json:"purevpn"`
	Surfshark    bool          `json:"surfshark"`
	Torguard     bool          `json:"torguard"`
	VPNUnlimited bool          `json:"vpnunlimited"`
	Vyprvpn      bool          `json:"vyprvpn"`
	Wevpn        bool          `json:"wevpn"`
	Windscribe   bool          `json:"windscribe"`
	// The two below should be used in CLI mode only
	CLI bool `json:"-"`
}

func (settings *Updater) String() string {
	return strings.Join(settings.lines(), "\n")
}

func (settings *Updater) lines() (lines []string) {
	if settings.Period == 0 {
		return nil
	}

	lines = append(lines, lastIndent+"Updater:")

	lines = append(lines, indent+lastIndent+"Period: every "+settings.Period.String())

	return lines
}

func (settings *Updater) EnableAll() {
	settings.Cyberghost = true
	settings.HideMyAss = true
	settings.Ipvanish = true
	settings.Ivpn = true
	settings.Mullvad = true
	settings.Nordvpn = true
	settings.Privado = true
	settings.PIA = true
	settings.Privado = true
	settings.Privatevpn = true
	settings.Protonvpn = true
	settings.Purevpn = true
	settings.Surfshark = true
	settings.Torguard = true
	settings.VPNUnlimited = true
	settings.Vyprvpn = true
	settings.Wevpn = true
	settings.Windscribe = true
}

func (settings *Updater) read(r reader) (err error) {
	settings.EnableAll()
	// use cloudflare in plaintext to not be blocked by DNS over TLS by default.
	// If a plaintext address is set in the DNS settings, this one will be used.
	// TODO use custom future encrypted DNS written in Go without blocking
	// as it's too much trouble to start another parallel unbound instance for now.
	settings.DNSAddress = "1.1.1.1"

	settings.Period, err = r.env.Duration("UPDATER_PERIOD", params.Default("0"))
	if err != nil {
		return fmt.Errorf("environment variable UPDATER_PERIOD: %w", err)
	}

	return nil
}
