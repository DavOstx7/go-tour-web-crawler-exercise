package main

import (
	"fmt"
	"sync"
)

type FetchRequest struct {
	url   string
	depth int
}

func fetchURL(url string, depth int, cacher *Cacher, fetcher Fetcher, requestChan chan<- *FetchRequest, wg *sync.WaitGroup) {
	defer wg.Done()

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

	for _, u := range urls {
		wg.Add(1)
		requestChan <- &FetchRequest{url: u, depth: depth - 1}
	}
}

func Crawl(url string, depth int, fetcher Fetcher) {
	requestChan := make(chan *FetchRequest)
	cacher := &Cacher{keys: make(map[string]bool)}
	var wg sync.WaitGroup

	go func() {
		for request := range requestChan {
			go fetchURL(request.url, request.depth, cacher, fetcher, requestChan, &wg)
		}
	}()

	wg.Add(1)
	requestChan <- &FetchRequest{url: url, depth: depth}

	wg.Wait()
	close(requestChan)
}

func main() {
	Crawl("https://golang.org/", 4, fetcher)
}
