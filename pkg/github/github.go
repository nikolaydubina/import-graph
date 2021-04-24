package github

import (
	"net/url"
	"strings"
)

func ParseGitHubURL(repoURL url.URL) (owner, repoName string) {
	parts := []string{}
	// Filtering out empty strings
	for _, p := range strings.Split(repoURL.EscapedPath(), "/") {
		if p != "" {
			parts = append(parts, p)
		}
	}
	if len(parts) != 2 {
		return "", ""
	}
	return parts[0], parts[1]
}
