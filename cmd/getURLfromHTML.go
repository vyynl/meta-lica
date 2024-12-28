package cmd

import (
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

func GetUrlFromHtml(htmlBody, rawBaseURL string) ([]string, error) {
	// Parsing rawBaseURL so that it can be used later to resolve relative references
	baseUrl, err := url.Parse(rawBaseURL)
	if err != nil {
		return []string{}, err
	}

	// Creating a reader object that is then parsed into a tree of html.Node objects
	reader := strings.NewReader(htmlBody)
	doc, err := html.Parse(reader)
	if err != nil {
		return []string{}, err
	}

	// Traversing the tree and saving any anchored URLs from each node
	var output []string
	for n := range doc.Descendants() {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, j := range n.Attr {
				if j.Key == "href" {

					url, err := url.Parse(j.Val)
					if err != nil {
						return []string{}, err
					}

					if url.Host == "" {
						output = append(output, baseUrl.ResolveReference(url).String())
					} else {
						output = append(output, j.Val)
					}
				}
			}
		}
	}
	return output, nil
}
