package cmd

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func GetHtml(rawURL string) (string, error) {
	// HTTP request parsing
	resp, err := http.Get(rawURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Checking to confirm that resp has the valid html Body string to exp
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("unable to complete request: %s", resp.Status)
	}
	contentType := resp.Header.Get("Content-Type")
	if strings.Contains(contentType, "text/html") != true {
		return "", fmt.Errorf("invalid content-type: %s", contentType)
	}

	output, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(output), nil
}
