package cmd

import (
	"fmt"
	"net/url"
)

func (cfg *config) CrawlPage(rawCurrentURL string) {
	cfg.ConcurrencyControl <- struct{}{}
	defer func() {
		<-cfg.ConcurrencyControl
		cfg.Wg.Done()
	}()

	if cfg.checkPagesCap() == true {
		return
	}

	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Printf("Error - unable to parse url '%s': %v\n", rawCurrentURL, err)
		return
	}

	/* Skipping any URLs that don't share the same host as cfg.baseURL or
	   have already been visited */
	if currentURL.Hostname() != cfg.BaseURL.Hostname() {
		return
	}

	normCurrentURL, err := NormalizeUrl(rawCurrentURL)
	if err != nil {
		fmt.Printf("Error - NormalizeURL: %v\n", err)
		return
	}

	isFirst := cfg.addPageVist(normCurrentURL)
	if !isFirst {
		return
	}

	/* Crawling through the new URLs that we find so we can extract any
	   other pages that it may have */
	fmt.Printf("crawling: %s\n", rawCurrentURL)

	htmlBody, err := GetHtml(rawCurrentURL)
	if err != nil {
		fmt.Printf("Error - GetHtml: %v\n", err)
		return
	}

	extractedURLs, err := GetUrlFromHtml(htmlBody, cfg.BaseURL.String())
	if err != nil {
		fmt.Printf("Error - GetUrlFromHtml: %v\n", err)
		return
	}
	for _, next := range extractedURLs {
		cfg.Wg.Add(1)
		go cfg.CrawlPage(next)
	}
}
