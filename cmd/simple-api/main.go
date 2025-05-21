package main

import (
	"log"
	"net/http"
	"sync"

	"kan.com/round-robin-api/internal/config"
	seedsimulator "kan.com/round-robin-api/internal/seed-simulator"
	simpleapi "kan.com/round-robin-api/internal/simple-api"
)

func main() {
	var wg sync.WaitGroup
	portsList := config.Get().ApiPorts

	for _, port := range portsList {
		wg.Add(1)
		currentPort := port

		go func() {
			defer wg.Done()

			mux := http.NewServeMux()
			mux.HandleFunc("POST /", func(w http.ResponseWriter, r *http.Request) {
				simpleapi.PostRoot(w, r, seedsimulator.New())
			})

			log.Printf("Starting server on %s\n", currentPort)
			server := &http.Server{
				Addr:    currentPort,
				Handler: mux,
			}

			if err := server.ListenAndServe(); err != nil {
				log.Fatalf("Error on port %s with error: %s\n", currentPort, err)
			}
		}()
	}

	wg.Wait()
}
