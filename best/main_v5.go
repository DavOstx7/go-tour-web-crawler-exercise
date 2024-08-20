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
	Process(r *CrawlResult)
}

type FakeProcessor struct{}

func (p *FakeProcessor) Process(r *CrawlResult) {
	if r.err != nil {
		fmt.Println(r.err)
		return
	}

	fmt.Printf("found: %s %q\n", r.url, r.body)
}

type Crawler struct {
	cacher    *Cacher
	fetcher   Fetcher
	processor Processor
	wg        *sync.WaitGroup
	requests  chan *CrawlRequest
	results   chan *CrawlResult
}

func (c *Crawler) sendRequest(request *CrawlRequest) {
	c.wg.Add(1)
	c.requests <- request
}

func (c *Crawler) waitAndClose() {
	c.wg.Wait()
	close(c.results)
	close(c.requests)
}

func (c *Crawler) listenForRequests() {
	for request := range c.requests {
		go c.crawl(request)
	}
}

func (c *Crawler) listenForResults() {
	for result := range c.results {
		c.processor.Process(result)

	}
}

func (c *Crawler) crawl(r *CrawlRequest) {
	defer c.wg.Done()

	if r.depth <= 0 || c.cacher.IsCached(r.url) {
		return
	}

	c.cacher.Cache(r.url)
	body, urls, err := c.fetcher.Fetch(r.url)
	c.results <- &CrawlResult{r.url, body, err}

	if err != nil {
		return
	}

	subDepth := r.depth - 1
	if subDepth == 0 {
		return
	}

	for _, subURL := range urls {
		if c.cacher.IsCached(subURL) {
			continue
		}
		c.sendRequest(&CrawlRequest{url: subURL, depth: subDepth})
	}
}

func (c *Crawler) Crawl(url string, depth int) {
	// Listen for requests
	go c.listenForRequests()

	// Send the first crawl request to the channel
	c.sendRequest(&CrawlRequest{url: url, depth: depth})

	// Wait for the wait group to finish, and then close the channel
	go c.waitAndClose()

	c.listenForResults()
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher) {
	crawler := Crawler{
		cacher:    &Cacher{keys: make(map[string]bool)},
		fetcher:   fetcher,
		processor: &FakeProcessor{},
		wg:        &sync.WaitGroup{},
		requests:  make(chan *CrawlRequest),
		results:   make(chan *CrawlResult),
	}

	// Listen for requests
	crawler.Crawl(url, depth)
}

func main() {
	Crawl("https://golang.org/", 4, fetcher)
}
