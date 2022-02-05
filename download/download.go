package download

import (
	"strings"
)

// Clean trims any URL parts, like the scheme or username, that might be present
// in a user-submitted URL
func Clean(path string) string {
	return trimUsername(trimScheme(path))
}

// trimScheme removes a scheme (e.g. https://) from the URL for more
// convenient pasting from browsers.
func trimScheme(repo string) string {
	schemeSep := "://"
	schemeSepIdx := strings.Index(repo, schemeSep)
	if schemeSepIdx > -1 {
		return repo[schemeSepIdx+len(schemeSep):]
	}

	return repo
}

// trimUsername removes the username for a URL, if it is present
func trimUsername(repo string) string {
	usernameSep := "@"
	usernameSepIdx := strings.Index(repo, usernameSep)
	if usernameSepIdx > -1 {
		return repo[usernameSepIdx+len(usernameSep):]
	}

	return repo
}
