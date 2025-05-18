package connectionpool

import (
	"errors"
	"testing"
	"time"
)

func TestGetUrl(t *testing.T) {
	cp := New([]string{":2222"})

	conn, _ := cp.GetConnection()

	if conn.GetUrl() != "http://localhost:2222" {
		t.Errorf("url should be http://localhost:2222 but is %s", conn.GetUrl())
	}
}

func TestAnalyse(t *testing.T) {
	runAsync = false

	cp := New([]string{":2222"})

	conn, _ := cp.GetConnection()

	conn.Analyse(nil, errors.New("some error"), time.Now())
}
