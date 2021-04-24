package awesomelists

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type AwesomeListsChecker struct {
	HTTPClient *http.Client
}

// IsMentioned returns if repo at github URL is mentioned in any of the lists
// Fetches lists from GitHub.
func (c *AwesomeListsChecker) IsMentioned(ghURL url.URL) (bool, error) {
	resp, err := http.Get("https://raw.githubusercontent.com/avelino/awesome-go/master/README.md")
	if err != nil {
		return false, fmt.Errorf("can not fetch awesome list go: %w", err)
	}
	defer func() { resp.Body.Close() }()

	str, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("can not read body response: %w", err)
	}

	return strings.Contains(string(str), ghURL.String()), nil
}
