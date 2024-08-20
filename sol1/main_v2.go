package main

import (
	"fmt"
	"sync"
)

type CrawlResult struct {
	url  string
	body string
}

func crawlURLs(url string, depth int, ch chan<- *CrawlResult, cacher *Cacher, fetcher Fetcher) {
	if depth == 0 || cacher.IsCached(url) {
		return
	}

	cacher.Cache(url)
	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	ch <- &CrawlResult{url: url, body: body}

	var wg sync.WaitGroup
	for _, u := range urls {
		if cacher.IsCached(u) {
			continue
		}

		wg.Add(1)
		go func(u string) {
			crawlURLs(u, depth-1, ch, cacher, fetcher)
			wg.Done()
		}(u)
	}
	wg.Wait()
}

func Crawl(url string, depth int, fetcher Fetcher) {
	ch := make(chan *CrawlResult)
	cacher := &Cacher{keys: make(map[string]bool)}

	go func() {
		crawlURLs(url, depth, ch, cacher, fetcher)
		close(ch)
	}()

	for result := range ch {
		fmt.Printf("found: %s %q\n", result.url, result.body)
	}
}

func main() {
	Crawl("https://golang.org/", 4, fetcher)
}
