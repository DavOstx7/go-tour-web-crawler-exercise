package main

import (
	"fmt"
)

type CrawlResult struct {
	url  string
	body string
}

func crawlURLs(url string, depth int, ch chan<- *CrawlResult, cacher *Cacher, fetcher Fetcher) {
	if depth == 0 || cacher.IsCached(url) {
		return
	}

	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	ch <- &CrawlResult{url: url, body: body}
	cacher.Cache(url)

	for _, u := range urls {
		crawlURLs(u, depth-1, ch, cacher, fetcher)
	}
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
