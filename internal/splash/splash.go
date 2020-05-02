package splash

import (
	"fmt"
	"strings"
	"time"

	"github.com/kyokomi/emoji"
	"github.com/qdm12/private-internet-access-docker/internal/constants"
)

// Splash returns the welcome spash message
func Splash(version, vcsRef, buildDate string) string {
	lines := title()
	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf("Running version %s built on %s (commit %s)", version, buildDate, vcsRef))
	lines = append(lines, "")
	lines = append(lines, announcement()...)
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

func announcement() []string {
	if len(constants.Announcement) == 0 {
		return nil
	}
	expirationDate, _ := time.Parse("2006-01-02", constants.AnnouncementExpiration) // error covered by a unit test
	if time.Now().After(expirationDate) {
		return nil
	}
	return []string{emoji.Sprint(":mega: ") + constants.Announcement}
}

func links() []string {
	return []string{
		emoji.Sprint(":wrench: ") + "Need help? " + constants.IssueLink,
		emoji.Sprint(":computer: ") + "Email? quentin.mcgaw@gmail.com",
		emoji.Sprint(":coffee: ") + "Slack? Join from the Slack button on Github",
		emoji.Sprint(":money_with_wings: ") + "Help me? https://github.com/sponsors/qdm12",
	}
}
