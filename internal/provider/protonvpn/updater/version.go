package updater

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

// getMostRecentStableTag finds the latest proton-vpn-gtk-app semver tag,
// in order to use it in the x-pm-appversion http request header. Because if we do
// fall behind on versioning, Proton doesn't like it because they like to create
// complications where there is no need for it. Hence this function.
func getMostRecentStableTag(ctx context.Context, client *http.Client) (version string, err error) {
	const url = "https://api.github.com/repos/ProtonVPN/proton-vpn-gtk-app/tags?per_page=30"

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("creating request: %w", err)
	}
	request.Header.Set("Accept", "application/vnd.github.v3+json")

	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%w: %s", ErrHTTPStatusCodeNotOK, response.Status)
	}

	decoder := json.NewDecoder(response.Body)
	var data []struct {
		Name string `json:"name"`
	}
	err = decoder.Decode(&data)
	if err != nil {
		return "", fmt.Errorf("decoding JSON response: %w", err)
	}

	// Find the first valid semver tag (tags are returned newest first by GitHub)
	regexVersion := regexp.MustCompile(`^v(\d+\.\d+\.\d+)$`)
	for _, tag := range data {
		if m := regexVersion.FindStringSubmatch(tag.Name); m != nil {
			version = "linux-vpn@" + strings.TrimPrefix(tag.Name, "v")
			return version, nil
		}
	}

	return "", fmt.Errorf("no valid semver tag found in proton-vpn-gtk-app tags")
}
