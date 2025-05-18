package connectionpool

import (
	"net/http"
	"net/url"
	"time"
)

const (
	scheme   = "http"
	hostname = "localhost"
)

type connection struct {
	url    *url.URL
	health health
}

func (c *connection) GetUrl() string {
	return c.url.String()
}

func (c *connection) Analyse(resp *http.Response, err error, startTime time.Time) {
	latency := time.Since(startTime).Milliseconds()
	isErr := err != nil || resp.StatusCode != http.StatusOK
	c.health.analyse(isErr, latency)
}
