package simpleapi

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPostRoot(t *testing.T) {
	t.Run("Unsupported Content-Type", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(`{}`)))
		req.Header.Set("Content-Type", "text/plain")

		rr := httptest.NewRecorder()
		PostRoot(rr, req)

		resp := rr.Result()
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnsupportedMediaType {
			t.Errorf("expected status 415, got %d", resp.StatusCode)
		}
	})

	t.Run("Valid JSON Request", func(t *testing.T) {
		jsonData := []byte(`{"hello":"world"}`)
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		PostRoot(rr, req)

		resp := rr.Result()
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}

		body, _ := io.ReadAll(resp.Body)
		if !bytes.Equal(body, jsonData) {
			t.Errorf("expected body %s, got %s", jsonData, body)
		}

		if ct := resp.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("expected Content-Type 'application/json', got '%s'", ct)
		}
	})
}
