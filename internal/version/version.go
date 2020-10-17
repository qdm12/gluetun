package version

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/qdm12/gluetun/internal/logging"
)

// GetMessage returns a message for the user describing if there is a newer version
// available. It should only be called once the tunnel is established.
func GetMessage(ctx context.Context, version, commitShort string, client *http.Client) (message string, err error) {
	if version == "latest" {
		// Find # of commits between current commit and latest commit
		commitsSince, err := getCommitsSince(ctx, client, commitShort)
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
	tagName, name, releaseTime, err := getLatestRelease(ctx, client)
	if err != nil {
		return "", fmt.Errorf("cannot get version information: %w", err)
	}
	if tagName == version {
		return fmt.Sprintf("You are running the latest release %s", version), nil
	}
	timeSinceRelease := logging.FormatDuration(time.Since(releaseTime))
	return fmt.Sprintf("There is a new release %s (%s) created %s ago",
			tagName, name, timeSinceRelease),
		nil
}

func getLatestRelease(ctx context.Context, client *http.Client) (tagName, name string, time time.Time, err error) {
	releases, err := getGithubReleases(ctx, client)
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

func getCommitsSince(ctx context.Context, client *http.Client, commitShort string) (n int, err error) {
	commits, err := getGithubCommits(ctx, client)
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
