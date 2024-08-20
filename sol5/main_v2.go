package main

import (
	"fmt"
	"sync"
)

func Crawl(url string, depth int, cacher *Cacher, fetcher Fetcher, wg *sync.WaitGroup) {
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
	for _, u := range urls {
		if !cacher.IsCached(u) {
			wg.Add(1)
			go Crawl(u, depth-1, cacher, fetcher, wg)
		}
	}

	return
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	cacher := &Cacher{keys: make(map[string]bool)}
	Crawl("https://golang.org/", 4, cacher, fetcher, &wg)

	wg.Wait()
}
