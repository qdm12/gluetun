package splash

import (
	"strings"
	"time"

	"github.com/kyokomi/emoji"
	"github.com/qdm12/private-internet-access-docker/internal/constants"
)

func Splash() string {
	lines := title()
	lines = append(lines, annoucement()...)
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
		return []string{emoji.Sprint(":rotating_light: ") + constants.Annoucement}
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
