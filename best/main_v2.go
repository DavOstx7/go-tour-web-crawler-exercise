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

func crawl(url string, depth int, ch chan<- *CrawlResult, cacher *Cacher, fetcher Fetcher, wg *sync.WaitGroup) {
	defer wg.Done()

	if depth <= 0 {
		return
	}

	cacher.Cache(url)
	body, urls, err := fetcher.Fetch(url)
	ch <- &CrawlResult{url: url, body: body, err: err}

	if len(urls) == 0 || depth == 1 {
		return
	}

	for _, u := range urls {
		if cacher.IsCached(u) {
			continue
		}

		wg.Add(1)
		go crawl(u, depth-1, ch, cacher, fetcher, wg)
	}
}

func Crawl(url string, depth int, fetcher Fetcher) {
	ch := make(chan *CrawlResult)
	cacher := &Cacher{keys: make(map[string]bool)}
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		crawl(url, depth, ch, cacher, fetcher, &wg)
		wg.Wait()
		close(ch)
	}()

	for result := range ch {
		if result.err != nil {
			fmt.Println(result.err)
			continue
		}
		fmt.Printf("found: %s %q\n", result.url, result.body)
	}
}

func main() {
	Crawl("https://golang.org/", 4, fetcher)
}
