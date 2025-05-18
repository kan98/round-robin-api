package main

import (
	"log"
	"net/http"
	"sync"

	"kan.com/round-robin-api/internal/config"
	simpleapi "kan.com/round-robin-api/internal/simple-api"
)

func main() {
	http.HandleFunc("/", simpleapi.PostRoot)
	var wg sync.WaitGroup
	portsList := config.Get().ApiPorts

	for _, v := range portsList {
		wg.Add(1)
		go func() {
			defer wg.Done()
			log.Printf("Starting server on %s\n", v)
			if err := http.ListenAndServe(v, nil); err != nil {
				log.Fatalf("Error on port %s with error: %s\n", v, err)
			}
		}()
	}

	wg.Wait()
}
