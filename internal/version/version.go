package version

import (
	"fmt"
	"net/http"
	"time"
)

// GetMessage returns a message for the user describing if there is a newer version
// available. It should only be called once the tunnel is established.
func GetMessage(version, commitShort string, client *http.Client) (message string, err error) {
	if version == "latest" {
		// Find # of commits between current commit and latest commit
		commitsSince, err := getCommitsSince(client, commitShort)
		if err != nil {
			return "", fmt.Errorf("cannot get version information: %w", err)
		} else if commitsSince == 0 {
			return fmt.Sprintf("You are running on the bleeding edge of %s!", version), nil
		}
		commits := "commits"
		if commitsSince == 1 {
			commits = "commit"
		}
		return fmt.Sprintf("You are running %d %s behind the most recent %s", commitsSince, commits, version), nil
	}
	tagName, name, releaseTime, err := getLatestRelease(client)
	if err != nil {
		return "", fmt.Errorf("cannot get version information: %w", err)
	}
	if tagName == version {
		return fmt.Sprintf("You are running the latest release %s", version), nil
	}
	timeSinceRelease := formatDuration(time.Since(releaseTime))
	return fmt.Sprintf("There is a new release %s (%s) created %s ago",
			tagName, name, timeSinceRelease),
		nil
}

func formatDuration(duration time.Duration) string {
	switch {
	case duration < time.Minute:
		seconds := int(duration.Round(time.Second).Seconds())
		if seconds < 2 {
			return fmt.Sprintf("%d second", seconds)
		}
		return fmt.Sprintf("%d seconds", seconds)
	case duration <= time.Hour:
		minutes := int(duration.Round(time.Minute).Minutes())
		if minutes == 1 {
			return "1 minute"
		}
		return fmt.Sprintf("%d minutes", minutes)
	case duration < 48*time.Hour:
		hours := int(duration.Truncate(time.Hour).Hours())
		return fmt.Sprintf("%d hours", hours)
	default:
		days := int(duration.Truncate(time.Hour).Hours() / 24)
		return fmt.Sprintf("%d days", days)
	}
}

func getLatestRelease(client *http.Client) (tagName, name string, time time.Time, err error) {
	releases, err := getGithubReleases(client)
	if err != nil {
		return "", "", time, err
	}
	for _, release := range releases {
		if release.Prerelease {
			continue
		}
		return release.TagName, release.Name, release.PublishedAt, nil
	}
	return "", "", time, fmt.Errorf("no releases found")
}

func getCommitsSince(client *http.Client, commitShort string) (n int, err error) {
	commits, err := getGithubCommits(client)
	if err != nil {
		return 0, err
	}
	for i := range commits {
		if commits[i].Sha[:7] == commitShort {
			return n, nil
		}
		n++
	}
	return 0, fmt.Errorf("no commit matching %q was found", commitShort)
}
