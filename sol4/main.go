package main

import (
	"fmt"
	"sync"
)

type CrawlResult struct {
	url   string
	depth int

	body string
	urls []string
	err  error
}

func crawlURL(url string, depth int, fetcher Fetcher, ch chan<- *CrawlResult) {
	body, urls, err := fetcher.Fetch(url)
	ch <- &CrawlResult{url, depth, body, urls, err}
}

func Crawl(url string, depth int, fetcher Fetcher) {
	ch := make(chan *CrawlResult)
	cacher := Cacher{keys: map[string]bool{url: true}}

	var wg sync.WaitGroup
	wg.Add(1)
	go crawlURL(url, depth, fetcher, ch)

	go func() {
		for r := range ch {
			if r.err != nil {
				fmt.Println(r.err)
			}
			fmt.Printf("found: %s %q\n", r.url, r.body)

			if r.depth > 0 {
				for _, u := range r.urls {
					if cacher.IsCached(u) {
						continue
					}
					cacher.Cache(u)

					wg.Add(1)
					go crawlURL(u, r.depth-1, fetcher, ch)
				}
			}
			wg.Done()
		}
	}()

	wg.Wait()
	close(ch)
}

func main() {
	Crawl("https://golang.org/", 4, fetcher)
}
