package main

import (
	"fmt"
	"sync"
)

type FetchResult struct {
	url  string
	body string
	err  error
}

type CrawlRequest struct {
	url   string
	depth int
}

type Crawler struct {
	depth    int
	cacher   *Cacher
	fetcher  Fetcher
	wg       *sync.WaitGroup
	requests chan CrawlRequest
	results  chan FetchResult
}

func (c *Crawler) Crawl(url string, depth int) {
	defer c.wg.Done()

	if depth <= 0 {
		return
	}

	body, urls, err := c.fetcher.Fetch(url)
	c.results <- FetchResult{url, body, err}

	if len(urls) == 0 {
		return
	}

	c.wg.Add(len(urls))
	for _, url := range urls {
		c.requests <- CrawlRequest{url, depth - 1}
	}
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher) (results chan FetchResult) {
	results = make(chan FetchResult)
	cacher := &Cacher{keys: make(map[string]bool)}

	crawler := Crawler{
		requests: make(chan CrawlRequest),
		depth:    depth,
		cacher:   cacher,
		fetcher:  fetcher,
		results:  results,
		wg:       &sync.WaitGroup{},
	}

	crawler.wg.Add(1)

	// Listen for requests, pass them through to the caller if they aren't duplicates.
	go func() {
		for request := range crawler.requests {
			if crawler.cacher.IsCached(request.url) {
				crawler.wg.Done()
				continue
			}

			crawler.cacher.Cache(request.url)
			go crawler.Crawl(request.url, request.depth)
		}
	}()

	// Wait for the wait group to finish, and then close the channel
	go func() {
		crawler.wg.Wait()
		close(results)
	}()

	// Send the first crawl request to the channel
	crawler.requests <- CrawlRequest{url, depth}

	return
}

func main() {
	results := Crawl("https://golang.org/", 4, fetcher)
	for result := range results {
		if result.err != nil {
			fmt.Println(result.err)
			continue
		}
		fmt.Printf("found: %s %q\n", result.url, result.body)
	}

}
