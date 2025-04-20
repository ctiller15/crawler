package internal

import (
	"net/url"
	"strings"
)

func normalizeURL(rawURL string) (string, error) {
	result, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	trimmedPath, _ := strings.CutSuffix(result.Path, "/")

	return strings.ToLower(result.Host + trimmedPath), nil
}
