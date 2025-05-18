package loadbalancerapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	connectionpool "kan.com/round-robin-api/internal/connection-pool"
)

var dummyPool = &connectionpool.MockConnectionPool{}

func TestPostRoot(t *testing.T) {
	// Create a new HTTP request with the desired URL and method
	req, err := http.NewRequest("POST", "http://localhost:2222/", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Set the Content-Type header to application/json
	req.Header.Set("Content-Type", "application/json")

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		PostRoot(w, r, dummyPool)
	})
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the Content-Type header in the response
	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Handler returned wrong Content-Type header: got %v want %v", contentType, "application/json")
	}
}
