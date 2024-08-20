package main

import (
	"fmt"
	"sync"
)

type FetchResult struct {
	url  string
	body string
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

	c.cacher.Cache(url)
	body, urls, err := c.fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	c.results <- FetchResult{url, body}

	if depth == 1 {
		return
	}

	for _, u := range urls {
		if c.cacher.IsCached(u) {
			continue
		}
		c.wg.Add(1)
		c.requests <- CrawlRequest{u, depth - 1}
	}
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher) (results chan FetchResult) {
	if depth <= 0 {
		return nil
	}

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

	// Listen for requests
	go func() {
		for request := range crawler.requests {
			go crawler.Crawl(request.url, request.depth)
		}
	}()

	crawler.wg.Add(1)

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
		fmt.Printf("found: %s %q\n", result.url, result.body)
	}

}
