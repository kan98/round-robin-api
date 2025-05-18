package loadbalancerapi

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	connectionpool "kan.com/round-robin-api/internal/connection-pool"
)

var client = &http.Client{
	Timeout: 10 * time.Second,
}

func PostRoot(w http.ResponseWriter, r *http.Request, pool connectionpool.ConnectionPool) {
	conn, err := pool.GetConnection()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	startTime := time.Now()
	resp, err := forwardRequest(r, conn.GetUrl())
	conn.Analyse(resp, err, startTime)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	io.Copy(w, resp.Body)
}

func forwardRequest(r *http.Request, url string) (*http.Response, error) {
	req, err := http.NewRequest(r.Method, url, r.Body)
	req.Header = r.Header
	if err != nil {
		log.Fatalf("Error creating new request: %v", err)
		return nil, err
	}

	fmt.Println("using connection:", url)
	return client.Do(req)
}
