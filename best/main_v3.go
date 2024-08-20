package main

import (
	"fmt"
	"sync"
)

type CrawlResult struct {
	url  string
	body string
	err  error
}

type CrawlRequest struct {
	url   string
	depth int
}

type Crawler struct {
	cacher  *Cacher
	fetcher Fetcher
	wg      *sync.WaitGroup
}

func (c *Crawler) SendRequest(url string, depth int, requests chan<- *CrawlRequest) {
	c.wg.Add(1)
	requests <- &CrawlRequest{url: url, depth: depth}
}

func (c *Crawler) Crawl(r *CrawlRequest, requests chan<- *CrawlRequest, results chan<- *CrawlResult) {
	defer c.wg.Done()

	if r.depth <= 0 {
		return
	}

	c.cacher.Cache(r.url)
	body, urls, err := c.fetcher.Fetch(r.url)
	results <- &CrawlResult{r.url, body, err}

	if err != nil || r.depth == 1 {
		return
	}

	for _, u := range urls {
		if c.cacher.IsCached(u) {
			continue
		}
		c.SendRequest(u, r.depth-1, requests)
	}
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher) chan *CrawlResult {
	requests := make(chan *CrawlRequest)
	results := make(chan *CrawlResult)

	crawler := Crawler{
		cacher:  &Cacher{keys: make(map[string]bool)},
		fetcher: fetcher,
		wg:      &sync.WaitGroup{},
	}

	// Listen for requests
	go func() {
		for r := range requests {
			go crawler.Crawl(r, requests, results)
		}
	}()

	// Send the first crawl request to the channel
	crawler.SendRequest(url, depth, requests)

	// Wait for the wait group to finish, and then close the channel
	go func() {
		crawler.wg.Wait()
		close(results)
		close(requests)
	}()

	for r := range results {
		if r.err != nil {
			fmt.Println(r.err)
			continue
		}
		fmt.Printf("found: %s %q\n", r.url, r.body)
	}
}

func main() {
	Crawl("https://golang.org/", 4, fetcher)
}
