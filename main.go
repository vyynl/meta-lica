package main

import (
	"fmt"
	"os"
	"strconv"
	"vyynl/meta-lica/cmd"
)

func main() {
	// Checking that user input is properly formatted
	if len(os.Args) < 2 {
		fmt.Println("no website provided\nusage: 'rawBaseURL' maxConcurrency maxPages")
		os.Exit(1)
	}
	if len(os.Args) > 4 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	rawBaseURL := os.Args[1]
	maxConcurrency, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("improper maxConcurrency formatting: must be integer")
	}
	maxPages, err := strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Println("improper maxPages formatting: must be integer")
	}

	cfg, err := cmd.Configure(rawBaseURL, maxConcurrency, maxPages)
	if err != nil {
		fmt.Printf("Error - Configure: %v", err)
		os.Exit(1)
	}

	fmt.Printf("starting crawl of: %s...\n", rawBaseURL)

	cfg.Wg.Add(1)
	go cfg.CrawlPage(rawBaseURL)
	cfg.Wg.Wait()

	for normalizedURL, count := range cfg.Pages {
		fmt.Printf("%d - %s\n", count, normalizedURL)
	}

	cfg.PrintReport()
}
