package updater

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"

	"golang.org/x/mod/semver"
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

	// Sort tags by semver. Invalid tags are placed at the end and we ignore them.
	// Yes, proton does push invalid semver tag names sometimes. Good job yet again.
	sort.Slice(data, func(i, j int) bool {
		return semver.Compare(data[i].Name, data[j].Name) > 0
	})

	version = "linux-vpn@" + data[0].Name[1:] // remove leading v
	return version, nil
}
