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

func (settings *PublicIP) String() string {
	return strings.Join(settings.lines(), "\n")
}

func (settings *PublicIP) lines() (lines []string) {
	if settings.Period == 0 {
		lines = append(lines, lastIndent+"Public IP getter: disabled")
		return lines
	}

	lines = append(lines, lastIndent+"Public IP getter:")
	lines = append(lines, indent+lastIndent+"Fetch period: "+settings.Period.String())
	lines = append(lines, indent+lastIndent+"IP file: "+string(settings.IPFilepath))

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
	if err != nil {
		return err
	}
	settings.IPFilepath = models.Filepath(filepathStr)

	return nil
}
