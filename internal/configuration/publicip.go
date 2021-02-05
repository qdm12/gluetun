package configuration

import (
	"strings"
	"time"

	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/golibs/params"
)

type PublicIP struct {
	Period     time.Duration   `json:"period"`
	IPFilepath models.Filepath `json:"ip_filepath"`
}

func (p *PublicIP) String() string {
	return strings.Join(p.lines(), "\n")
}

func (p *PublicIP) lines() (lines []string) {
	if p.Period == 0 {
		lines = append(lines, lastIndent+"Public IP getter: disabled")
		return lines
	}

	lines = append(lines, lastIndent+"Public IP getter:")
	lines = append(lines, indent+lastIndent+"Fetch period: "+p.Period.String())
	lines = append(lines, indent+lastIndent+"IP file: "+string(p.IPFilepath))

	return lines
}

func (settings *PublicIP) read(r reader) (err error) {
	settings.Period, err = r.env.Duration("PUBLICIP_PERIOD", params.Default("12h"))
	if err != nil {
		return err
	}

	filepathStr, err := r.env.Path("PUBLICIP_FILE", params.CaseSensitiveValue(),
		params.Default("/tmp/gluetun/ip"),
		params.RetroKeys([]string{"IP_STATUS_FILE"}, r.onRetroActive))
	settings.IPFilepath = models.Filepath(filepathStr)

	return nil
}
