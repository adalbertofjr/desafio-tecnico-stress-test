package main

import (
	"flag"
	"fmt"
	"net/http"
	"sync"
)

func main() {
	url := flag.String("url", "", "The URL to stress test")
	requests := flag.Int("requests", 0, "Number of requests to perform")
	concurrency := flag.Int("concurrency", 0, "Number of multiple requests to make at a time")

	if *url == "" || *requests <= 0 || *concurrency <= 0 {
		fmt.Println("Usage: stresstest -url <URL> -requests <number> -concurrency <number>")
		return
	}

	flag.Parse()

	execute(*url, *requests, *concurrency)
}

func execute(url string, requests int, concurrency int) {
	mu := sync.Mutex{}
	wg := sync.WaitGroup{}
	counter := 0
	ch := make(chan struct{}, concurrency)

	for range requests {
		wg.Add(1)
		ch <- struct{}{}
		go func() {
			defer func() { <-ch }()
			makeRequest(&mu, &wg, url, &counter)
		}()
	}
	wg.Wait()
}

func makeRequest(mu *sync.Mutex, wg *sync.WaitGroup, url string, counter *int) {
	defer wg.Done()
	client := http.Client{}

	resp, err := client.Get(url)

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// body, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	panic(err)
	// }

	fmt.Printf("Request: %d\n", incrementCount(mu, counter))
}

func incrementCount(mu *sync.Mutex, count *int) int {
	mu.Lock()
	defer mu.Unlock()

	*count++

	return *count
}
