package splash

import (
	"fmt"
	"strings"
	"time"

	"github.com/kyokomi/emoji"
	"github.com/qdm12/private-internet-access-docker/internal/constants"
	"github.com/qdm12/private-internet-access-docker/internal/params"
)

// Splash returns the welcome spash message
func Splash(paramsReader params.ParamsReader) string {
	version := paramsReader.GetVersion()
	vcsRef := paramsReader.GetVcsRef()
	buildDate := paramsReader.GetBuildDate()
	lines := title()
	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf("Running version %s built on %s (commit %s)", version, buildDate, vcsRef))
	lines = append(lines, "")
	lines = append(lines, annoucement()...)
	lines = append(lines, "")
	lines = append(lines, links()...)
	return strings.Join(lines, "\n")
}

func title() []string {
	return []string{
		"=========================================",
		"============= PIA container =============",
		"========== An exquisite mix of ==========",
		"==== OpenVPN, Unbound, DNS over TLS, ====",
		"===== Shadowsocks, Tinyproxy and Go =====",
		"=========================================",
		"=== Made with " + emoji.Sprint(":heart:") + " by github.com/qdm12 ====",
		"=========================================",
	}
}

func annoucement() []string {
	timestamp := time.Now().UnixNano() / 1000000000
	if timestamp < constants.AnnoucementExpiration {
		return []string{emoji.Sprint(":mega: ") + constants.Annoucement}
	}
	return nil
}

func links() []string {
	return []string{
		emoji.Sprint(":wrench: ") + "Need help? " + constants.IssueLink,
		emoji.Sprint(":computer: ") + "Email? quentin.mcgaw@gmail.com",
		emoji.Sprint(":coffee: ") + "Slack? Join from the Slack button on Github",
		emoji.Sprint(":money_with_wings: ") + "Help me? https://github.com/sponsors/qdm12",
	}
}
