package main

import (
	"fmt"
)

func Crawl(url string, depth int, cacher *Cacher, fetcher Fetcher, ret chan string) {
	defer close(ret)

	if depth <= 0 || cacher.IsCached(url) {
		return
	}

	cacher.Cache(url)
	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		ret <- err.Error()
		return
	}

	ret <- fmt.Sprintf("found: %s %q", url, body)

	result := make([]chan string, len(urls))
	for i, u := range urls {
		result[i] = make(chan string)
		go Crawl(u, depth-1, cacher, fetcher, result[i])
	}

	for i := range result {
		for s := range result[i] {
			ret <- s
		}
	}
}

func main() {
	result := make(chan string)
	cacher := Cacher{keys: make(map[string]bool)}

	go Crawl("https://golang.org/", 4, &cacher, fetcher, result)

	for s := range result {
		fmt.Println(s)
	}
}
