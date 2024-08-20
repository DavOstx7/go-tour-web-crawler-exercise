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

func crawl(url string, depth int, cacher *Cacher, fetcher Fetcher, wg *sync.WaitGroup) {
	defer wg.Done()

	if depth <= 0 {
		return
	}

	cacher.Cache(url)
	body, urls, err := fetcher.Fetch(url)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("found: %s %q\n", url, body)
	if len(urls) == 0 || depth == 1 {
		return
	}

	for _, u := range urls {
		if cacher.IsCached(u) {
			continue
		}

		wg.Add(1)
		go crawl(u, depth-1, cacher, fetcher, wg)
	}
}

func Crawl(url string, depth int, fetcher Fetcher) {
	cacher := &Cacher{keys: make(map[string]bool)}
	var wg sync.WaitGroup

	wg.Add(1)
	crawl(url, depth, cacher, fetcher, &wg)
	wg.Wait()
}

func main() {
	Crawl("https://golang.org/", 4, fetcher)
}
