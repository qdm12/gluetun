package privateinternetaccess

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// replaceInErr is used to remove sensitive information from errors.
func ReplaceInErr(err error, substitutions map[string]string) error {
	s := ReplaceInString(err.Error(), substitutions)
	return errors.New(s) //nolint:goerr113
}

// replaceInString is used to remove sensitive information.
func ReplaceInString(s string, substitutions map[string]string) string {
	for old, new := range substitutions {
		s = strings.ReplaceAll(s, old, new)
	}
	return s
}

func makeNOKStatusError(response *http.Response, substitutions map[string]string) (err error) {
	url := response.Request.URL.String()
	url = ReplaceInString(url, substitutions)

	b, _ := io.ReadAll(response.Body)
	shortenMessage := string(b)
	shortenMessage = strings.ReplaceAll(shortenMessage, "\n", "")
	shortenMessage = strings.ReplaceAll(shortenMessage, "  ", " ")
	shortenMessage = ReplaceInString(shortenMessage, substitutions)

	return fmt.Errorf("%w: %s: %d %s: response received: %s",
		ErrHTTPStatusCodeNotOK, url, response.StatusCode,
		response.Status, shortenMessage)
}
