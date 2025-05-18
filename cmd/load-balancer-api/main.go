package main

import (
	"log"
	"net/http"

	"kan.com/round-robin-api/internal/config"
	connectionpool "kan.com/round-robin-api/internal/connection-pool"
	loadbalancerapi "kan.com/round-robin-api/internal/load-balancer-api"
)

func main() {
	connectionPool := connectionpool.New(config.Get().ApiPorts)
	http.HandleFunc("POST /", func(w http.ResponseWriter, r *http.Request) {
		loadbalancerapi.PostRoot(w, r, connectionPool)
	})
	loadbalancerPort := config.Get().LoadBalancerPort

	log.Printf("Starting server on %s\n", loadbalancerPort)
	if err := http.ListenAndServe(loadbalancerPort, nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
