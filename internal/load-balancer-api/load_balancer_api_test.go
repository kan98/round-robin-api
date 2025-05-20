package loadbalancerapi

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/h2non/gock"
	connectionpool "kan.com/round-robin-api/internal/connection-pool"
)

func TestPostRoot(t *testing.T) {
	req, _ := http.NewRequest("POST", "http://localhost:2222/", nil)
	req.Header.Set("Content-Type", "application/json")

	t.Run("successful test", func(t *testing.T) {
		gock.New("http://localhost:2222").
			Post("/").
			Reply(200).
			JSON(map[string]string{"foo": "bar"})

		var dummyPool = &connectionpool.MockConnectionPool{}

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			PostRoot(w, r, dummyPool)
		})
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("want status %d but got %d", http.StatusOK, rr.Code)
		}

		if rr.Header().Get("Content-Type") != "application/json" {
			t.Errorf("want Content-Type %q but got %q", "application/json", rr.Header().Get("Content-Type"))
		}
	})

	t.Run("status code is same as downstream", func(t *testing.T) {
		gock.New("http://localhost:2222").
			Post("/").
			Reply(429)

		var dummyPool = &connectionpool.MockConnectionPool{}

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			PostRoot(w, r, dummyPool)
		})
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusTooManyRequests {
			t.Errorf("want status %d but got %d", http.StatusTooManyRequests, rr.Code)
		}
	})

	t.Run("error from downstream api", func(t *testing.T) {
		gock.New("http://localhost:2222").
			Post("/").
			ReplyError(errors.New("some error"))

		var dummyPool = &connectionpool.MockConnectionPool{}

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			PostRoot(w, r, dummyPool)
		})
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusInternalServerError {
			t.Errorf("want status %d but got %d", http.StatusInternalServerError, rr.Code)
		}
	})

	t.Run("cannot get connection from pool", func(t *testing.T) {
		var dummyPool = &connectionpool.MockConnectionPool{
			GetConnectionToReturnErr: true,
		}

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			PostRoot(w, r, dummyPool)
		})
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusInternalServerError {
			t.Errorf("want status %d but got %d", http.StatusInternalServerError, rr.Code)
		}
	})
}
