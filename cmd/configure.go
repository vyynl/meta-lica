package cmd

import (
	"fmt"
	"net/url"
	"sync"
)

type config struct {
	Pages              map[string]int
	BaseURL            *url.URL
	Mu                 *sync.Mutex
	ConcurrencyControl chan struct{}
	PageControl        int
	Wg                 *sync.WaitGroup
}

func (cfg *config) addPageVist(normalizedURL string) (isFirst bool) {
	cfg.Mu.Lock()
	defer cfg.Mu.Unlock()

	if _, visited := cfg.Pages[normalizedURL]; visited {
		cfg.Pages[normalizedURL]++
		return false
	}

	cfg.Pages[normalizedURL] = 1
	return true
}

func (cfg *config) checkPagesCap() (atCap bool) {
	cfg.Mu.Lock()
	defer cfg.Mu.Unlock()

	return len(cfg.Pages) >= cfg.PageControl
}

type keyValuePair struct {
	url   string
	count int
}

func (cfg *config) sortPagesByCount() []keyValuePair {
	var sortSlice []keyValuePair

	for key, value := range cfg.Pages {
		sortSlice = append(sortSlice, keyValuePair{key, value})
	}
	return sortSlice
}

func (cfg *config) PrintReport() {
	fmt.Println()
	fmt.Println("=============================")
	fmt.Printf("  REPORT for %s\n", cfg.BaseURL.String())
	fmt.Println("=============================")
	fmt.Println()
	for _, i := range cfg.sortPagesByCount() {
		fmt.Printf("Found %d internal links to %s\n", i.count, i.url)
	}
}

func Configure(rawBaseURL string, maxConcurrency, maxPages int) (*config, error) {
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		return nil, fmt.Errorf("unable to parse base url: %w\n", err)
	}

	return &config{
		Pages:              make(map[string]int),
		BaseURL:            baseURL,
		Mu:                 &sync.Mutex{},
		ConcurrencyControl: make(chan struct{}, maxConcurrency),
		PageControl:        maxPages,
		Wg:                 &sync.WaitGroup{},
	}, nil
}
