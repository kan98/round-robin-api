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

	for i, port := range portsList {
		wg.Add(1)

		go func() {
			defer wg.Done()

			mux := http.NewServeMux()
			mux.HandleFunc("POST /", func(w http.ResponseWriter, r *http.Request) {
				simpleapi.PostRoot(w, r, seedsimulator.New(i))
			})

			log.Printf("Starting server on %s\n", port)
			server := &http.Server{
				Addr:    port,
				Handler: mux,
			}

			if err := server.ListenAndServe(); err != nil {
				log.Fatalf("Error on port %s with error: %s\n", port, err)
			}
		}()
	}

	wg.Wait()
}
