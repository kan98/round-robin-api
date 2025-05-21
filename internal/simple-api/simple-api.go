package simpleapi

import (
	"io"
	"net/http"
	"time"

	seedsimulator "kan.com/round-robin-api/internal/seed-simulator"
)

func PostRoot(w http.ResponseWriter, r *http.Request, seed *seedsimulator.Seed) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Unsupported Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	if seed.SleepTime() > 0 {
		time.Sleep(time.Duration(seed.SleepTime()) * time.Millisecond)
	}
	if seed.ToError() {
		http.Error(w, "simulated error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}
