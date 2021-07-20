// Package logging defines helper functions for logging.
package logging

import (
	"fmt"
	"strings"
	"time"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
)

// Splash returns the welcome spash message.
func Splash(buildInfo models.BuildInformation) string {
	lines := title()
	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf("Running version %s built on %s (commit %s)",
		buildInfo.Version, buildInfo.Created, buildInfo.Commit))
	lines = append(lines, "")
	lines = append(lines, announcement()...)
	lines = append(lines, "")
	lines = append(lines, links()...)
	return strings.Join(lines, "\n")
}

func title() []string {
	return []string{
		"=========================================",
		"================ Gluetun ================",
		"=========================================",
		"==== A mix of OpenVPN, DNS over TLS, ====",
		"======= Shadowsocks and HTTP proxy ======",
		"========= all glued up with Go ==========",
		"=========================================",
		"=========== For tunneling to ============",
		"======== your favorite VPN server =======",
		"=========================================",
		"=== Made with ❤️ by github.com/qdm12 ====",
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
	return []string{"📣" + constants.Announcement}
}

func links() []string {
	return []string{
		"🔧 Need help? " + constants.IssueLink,
		"💻 Email? quentin.mcgaw@gmail.com",
		"☕ Slack? Join from the Slack button on Github",
		"💰 Help me? https://github.com/sponsors/qdm12",
	}
}
