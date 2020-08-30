package version

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

type githubRelease struct {
	TagName     string    `json:"tag_name"`
	Name        string    `json:"name"`
	Prerelease  bool      `json:"prerelease"`
	PublishedAt time.Time `json:"published_at"`
}

type githubCommit struct {
	Sha    string `json:"sha"`
	Commit struct {
		Committer struct {
			Date time.Time `json:"date"`
		} `json:"committer"`
	}
}

func getGithubReleases(client *http.Client) (releases []githubRelease, err error) {
	const url = "https://api.github.com/repos/qdm12/gluetun/releases"
	response, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(b, &releases); err != nil {
		return nil, err
	}
	return releases, nil
}

func getGithubCommits(client *http.Client) (commits []githubCommit, err error) {
	const url = "https://api.github.com/repos/qdm12/gluetun/commits"
	response, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(b, &commits); err != nil {
		return nil, err
	}
	return commits, nil
}
