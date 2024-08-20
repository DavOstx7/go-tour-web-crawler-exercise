package main

import (
	"fmt"
)

func Crawl(url string, depth int, cacher *Cacher, fetcher Fetcher, exit chan bool) {
	// Fetch URLs in parallel.
	// Don't fetch the same URL twice.

	if depth <= 0 {
		exit <- true
		return
	}

	if cacher.IsCached(url) {
		exit <- true
		return
	}

	cacher.Cache(url)
	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		exit <- true
		return
	}
	fmt.Printf("found: %s %q\n", url, body)

	e := make(chan bool)
	for _, u := range urls {
		go Crawl(u, depth-1, cacher, fetcher, e)
	}

	// wait for all child gorountines to exit
	for i := 0; i < len(urls); i++ {
		<-e
	}
	exit <- true
}

func main() {
	exit := make(chan bool)
	cacher := &Cacher{keys: make(map[string]bool)}
	go Crawl("https://golang.org/", 4, cacher, fetcher, exit)
	<-exit
}