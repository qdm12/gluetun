package configuration

import (
	"strings"
	"time"

	"github.com/qdm12/golibs/params"
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

func (s *Updater) String() string {
	return strings.Join(s.lines(), "\n")
}

func (s *Updater) lines() (lines []string) {
	if s.Period == 0 {
		return nil
	}

	lines = append(lines, lastIndent+"Updater:")

	lines = append(lines, indent+lastIndent+"Period: every "+s.Period.String())

	return lines
}

func (settings *Updater) read(r reader) (err error) {
	settings.Cyberghost = true
	settings.Mullvad = true
	settings.Nordvpn = true
	settings.PIA = true
	settings.Purevpn = true
	settings.Surfshark = true
	settings.Vyprvpn = true
	settings.Windscribe = true
	settings.Stdout = false
	settings.CLI = false
	settings.DNSAddress = "127.0.0.1"

	settings.Period, err = r.env.Duration("UPDATER_PERIOD", params.Default("0"))
	if err != nil {
		return err
	}

	return nil
}
