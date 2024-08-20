package main

import (
	"fmt"
	"sync"
)

func crawlURLs(url string, depth int, cacher *Cacher, fetcher Fetcher) {
	if depth == 0 || cacher.IsCached(url) {
		return
	}

	cacher.Cache(url)
	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("found: %s %q\n", url, body)

	var wg sync.WaitGroup
	wg.Add(len(urls))
	for _, u := range urls {
		u := u
		go func() {
			crawlURLs(u, depth-1, cacher, fetcher)
			wg.Done()
		}()
	}
	wg.Wait()
}

func Crawl(url string, depth int, fetcher Fetcher) {
	cacher := &Cacher{keys: make(map[string]bool)}
	crawlURLs(url, depth, cacher, fetcher)
}

func main() {
	Crawl("https://golang.org/", 4, fetcher)
}
