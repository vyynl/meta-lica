package cmd

import (
	"net/url"
	"strings"
)

func NormalizeUrl(input string) (string, error) {
	u, err := url.Parse(input)
	if err != nil {
		return "", err
	}

	cleaned := strings.TrimRight(strings.ToLower(u.Hostname()+u.Path), "/")
	return cleaned, nil
}
