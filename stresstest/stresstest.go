package stresstest

import (
	"net/http"
	"sync"
	"time"
)

func Execute(url string, requests int, concurrency int) {
	startTime := time.Now()
	mu := sync.Mutex{}
	wg := sync.WaitGroup{}
	// counter := 0
	ch := make(chan struct{}, concurrency)

	var data Report
	data.url = url
	data.TotalRequests = requests
	data.SuccessRequests = 0
	data.FailureRequests = make(map[int]int)

	for range requests {
		wg.Add(1)
		ch <- struct{}{}
		go func() {
			defer func() { <-ch }()
			makeRequest(&mu, &wg, url, &data)
		}()
	}
	wg.Wait()
	data.TimeDuration = time.Since(startTime)
	PrintReport(data)
}

func makeRequest(mu *sync.Mutex, wg *sync.WaitGroup, url string, data *Report) {
	defer wg.Done()
	client := http.Client{}

	resp, err := client.Get(url)
	if err != nil {
		panic(err)
	}

	incrementCount(mu, &data.RequestsExecs)
	incrementaStatusCount(mu, resp.StatusCode, data)
}

func incrementCount(mu *sync.Mutex, count *int) int {
	mu.Lock()
	defer mu.Unlock()

	*count++

	return *count
}

func incrementaStatusCount(mu *sync.Mutex, statusCode int, data *Report) {
	mu.Lock()
	defer mu.Unlock()

	if statusCode >= 200 && statusCode < 300 {
		data.SuccessRequests++
	} else {
		data.FailureRequests[statusCode]++
	}
}
