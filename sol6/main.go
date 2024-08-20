package main

import (
	"fmt"
)

type FakeResult struct {
	body string
	url  string
	err  error
}

func crawlHelper(url string, depth int, cacher *Cacher, fetcher Fetcher, ch chan<- FakeResult, exit chan bool) {
	cacher.Cache(url)

	body, urls, err := fetcher.Fetch(url)
	ch <- FakeResult{body, url, err}

	if depth > 1 {
		children := 0

		for _, u := range urls {
			if !cacher.IsCached(u) { // only create a new goroutine if url hasn't been fetched
				children++
				go crawlHelper(u, depth-1, cacher, fetcher, ch, exit)
			}
		}

		for i := 0; i < children; i++ {
			<-exit // collect exit calls to prevent calling parent directly
		}
	}

	exit <- true // call parent when this goroutine and its children are done
}

func Crawl(url string, depth int, fetcher Fetcher) {
	cacher := Cacher{keys: map[string]bool{}}
	ch := make(chan FakeResult)
	exit := make(chan bool)

	go crawlHelper(url, depth, &cacher, fetcher, ch, exit)

	for { // infinite loop ...
		select { // ... with a listener
		case fakeResult := <-ch:
			if fakeResult.err != nil {
				fmt.Println(fakeResult.err)
			} else {
				fmt.Printf("found: %s %q\n", fakeResult.url, fakeResult.body)
			}
		case <-exit:
			return // terminate listening to all goroutines
		}
	}
}

func main() {
	Crawl("https://golang.org/", 4, fetcher)
}
