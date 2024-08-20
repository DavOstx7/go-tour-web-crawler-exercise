package main

import (
	"fmt"
	"sync"
)

type CrawlRequest struct {
	url   string
	depth int
}

type CrawlResult struct {
	url  string
	body string
	err  error
}

type Processor interface {
	Process(url string, body string, err error)
}

type FakeProcessor struct{}

func (p *FakeProcessor) Process(url string, body string, err error) {
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("found: %s %q\n", url, body)
}

type Crawler struct {
	cacher   *Cacher
	fetcher  Fetcher
	wg       *sync.WaitGroup
	requests chan *CrawlRequest
	results  chan *CrawlResult
}

func (c *Crawler) SendRequest(r *CrawlRequest) {
	c.wg.Add(1)
	c.requests <- r
}

func (c *Crawler) WaitAndClose() {
	c.wg.Wait()
	close(c.results)
	close(c.requests)
}

func (c *Crawler) CrawlRequests() {
	for request := range c.requests {
		go c.crawl(request.url, request.depth)
	}
}

func (c *Crawler) ProcessResults(p Processor) {
	for result := range c.results {
		p.Process(result.url, result.body, result.err)

	}
}

func (c *Crawler) crawl(url string, depth int) {
	defer c.wg.Done()

	if depth <= 0 || c.cacher.IsCached(url) {
		return
	}

	c.cacher.Cache(url)
	body, urls, err := c.fetcher.Fetch(url)
	c.results <- &CrawlResult{url, body, err}

	if err != nil {
		return
	}

	subDepth := depth - 1
	if subDepth == 0 {
		return
	}

	for _, subURL := range urls {
		if c.cacher.IsCached(subURL) {
			continue
		}
		c.SendRequest(&CrawlRequest{url: subURL, depth: subDepth})
	}
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher) {
	crawler := Crawler{
		cacher:   &Cacher{keys: make(map[string]bool)},
		fetcher:  fetcher,
		wg:       &sync.WaitGroup{},
		requests: make(chan *CrawlRequest),
		results:  make(chan *CrawlResult),
	}

	// Listen for requests
	go crawler.CrawlRequests()

	// Send the first crawl request to the channel
	crawler.SendRequest(&CrawlRequest{url: url, depth: depth})

	// Wait for the wait group to finish, and then close the channel
	go crawler.WaitAndClose()

	crawler.ProcessResults(&FakeProcessor{})
}

func main() {
	Crawl("https://golang.org/", 4, fetcher)
}
