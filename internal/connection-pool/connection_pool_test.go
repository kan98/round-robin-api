package connectionpool

import (
	"net/url"
	"testing"

	"kan.com/round-robin-api/internal/config"
)

func TestNew(t *testing.T) {
	t.Run("creates pool with correct conns", func(t *testing.T) {
		pool := New([]string{":2222", ":3333", ":4444"})

		expectedUrls := []string{
			"http://localhost:2222",
			"http://localhost:3333",
			"http://localhost:4444",
		}

		if len(pool.connections) != 3 {
			t.Error("expected 3 connections")
		}

		for i, conn := range pool.connections {
			if conn.GetUrl() != expectedUrls[i] {
				t.Error("conns incorrect or in the wrong order")
			}
		}
	})

	t.Run("creates pool with no conns when empty apiPorts env", func(t *testing.T) {
		pool := New([]string{})

		if len(pool.connections) != 0 {
			t.Error("expected 0 connections")
		}
	})
}

func TestGetConnection(t *testing.T) {
	t.Run("get healthy connection", func(t *testing.T) {
		pool := &connectionPool{
			connections: []connection{
				{
					url:    &url.URL{Scheme: "http", Host: "localhost:2222"},
					health: health{penaltyLeftTimes: 0},
				},
				{
					url:    &url.URL{Scheme: "http", Host: "localhost:3333"},
					health: health{penaltyLeftTimes: 1},
				},
			},
			currentIndex: 0,
		}

		conn, err := pool.GetConnection()
		if err != nil {
			t.Errorf("Expected no error, got %s", err)
		}
		if conn == nil {
			t.Error("No connection found")
		}
		if conn.GetUrl() != "http://localhost:2222" {
			t.Error("Wrong connection")
		}
	})

	t.Run("when not optimised, doesn't skip unhealthy connection", func(t *testing.T) {
		config.Reset()
		t.Setenv("optimiseConnPool", "false")

		pool := &connectionPool{
			connections: []connection{
				{
					url:    &url.URL{Scheme: "http", Host: "localhost:2222"},
					health: health{penaltyLeftTimes: 0},
				},
				{
					url:    &url.URL{Scheme: "http", Host: "localhost:3333"},
					health: health{penaltyLeftTimes: 1},
				},
			},
			currentIndex: 1,
		}

		conn, _ := pool.GetConnection()
		if conn.GetUrl() != "http://localhost:3333" {
			t.Error("Wrong connection")
		}
	})

	t.Run("when optimised, skips unhealthy connection", func(t *testing.T) {
		config.Reset()
		t.Setenv("optimiseConnPool", "true")

		pool := &connectionPool{
			connections: []connection{
				{
					url:    &url.URL{Scheme: "http", Host: "localhost:2222"},
					health: health{penaltyLeftTimes: 0},
				},
				{
					url:    &url.URL{Scheme: "http", Host: "localhost:3333"},
					health: health{penaltyLeftTimes: 1},
				},
			},
			currentIndex: 1,
		}

		conn, _ := pool.GetConnection()
		if conn.GetUrl() != "http://localhost:2222" {
			t.Error("Wrong connection")
		}
	})

	t.Run("when optimised and all unhealthy cons, just round robins next one", func(t *testing.T) {
		pool := &connectionPool{
			connections: []connection{
				{
					url:    &url.URL{Scheme: "http", Host: "localhost:2222"},
					health: health{penaltyLeftTimes: 5},
				},
				{
					url:    &url.URL{Scheme: "http", Host: "localhost:3333"},
					health: health{penaltyLeftTimes: 6},
				},
			},
			currentIndex: 1,
		}

		conn, _ := pool.GetConnection()
		if conn.GetUrl() != "http://localhost:3333" {
			t.Error("Wrong connection")
		}
	})

	t.Run("returns error for empty pool", func(t *testing.T) {
		pool := &connectionPool{
			connections: []connection{},
		}

		_, err := pool.GetConnection()
		if err == nil {
			t.Error("expected an error")
		}
	})
}
