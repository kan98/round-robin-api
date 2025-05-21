package main

import (
	"bytes"
	"fmt"
	"net/http"
	"sync"
	"time"
)

const (
	totalRequests = 500
	concurrency   = 100
	url           = "http://localhost:1111"
)

func worker(jobs <-chan int, isSuccess chan<- bool, wg *sync.WaitGroup) {
	defer wg.Done()
	client := &http.Client{}

	for range jobs {
		jsonData := []byte(`{"hello":"world"}`)
		resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonData))
		if err == nil && resp.StatusCode == http.StatusOK {
			isSuccess <- true
		} else {
			isSuccess <- false
		}
		if resp != nil {
			resp.Body.Close()
		}
	}
}

func main() {
	start := time.Now()

	jobs := make(chan int, totalRequests)
	isSuccess := make(chan bool, totalRequests)

	var wg sync.WaitGroup

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go worker(jobs, isSuccess, &wg)
	}

	for i := 0; i < totalRequests; i++ {
		jobs <- i
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(isSuccess)
	}()

	successes := 0
	for success := range isSuccess {
		if success {
			successes++
		}
	}

	fmt.Printf("Total Time: %.2f seconds\n", time.Since(start).Seconds())
	fmt.Printf("Total Successes: %d\n", successes)
	fmt.Printf("Total Errors: %d\n", totalRequests-successes)
}
