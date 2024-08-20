package main

import (
	"fmt"
	"sync"
)

type FetchRequest struct {
	url   string
	depth int
}

type FetchResult struct {
	url  string
	body string
}

func fetchURL(url string, depth int, cacher *Cacher, fetcher Fetcher, requestChan chan<- *FetchRequest, resultChan chan<- *FetchResult, wg *sync.WaitGroup) {
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

	resultChan <- &FetchResult{url: url, body: body}

	for _, u := range urls {
		if !cacher.IsCached(u) {
			wg.Add(1)
			requestChan <- &FetchRequest{url: u, depth: depth - 1}
		}
	}
}

func Crawl(url string, depth int, fetcher Fetcher) {
	requestChan := make(chan *FetchRequest)
	resultChan := make(chan *FetchResult)

	cacher := &Cacher{keys: make(map[string]bool)}
	var wg sync.WaitGroup

	go func() {
		for request := range requestChan {
			go fetchURL(request.url, request.depth, cacher, fetcher, requestChan, resultChan, &wg)
		}
	}()

	wg.Add(1)
	requestChan <- &FetchRequest{url: url, depth: depth}

	go func() {
		for result := range resultChan {
			fmt.Printf("found: %s %q\n", result.url, result.body)
		}
	}()

	wg.Wait()
	close(requestChan)
	close(resultChan)
}

func main() {
	Crawl("https://golang.org/", 4, fetcher)
}
