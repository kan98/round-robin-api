package connectionpool

import (
	"net/http"
	"time"
)

// ConnectionPool defines the interface for connection pool operations
type ConnectionPool interface {
	// GetConnection returns the next available connection from the pool
	GetConnection() (Connection, error)
}

// Connection defines the interface for a connection to a backend service
type Connection interface {
	// GetUrl returns the URL string for this connection
	GetUrl() string

	// Analyse evaluates the health of the connection based on response
	Analyse(resp *http.Response, err error, startTime time.Time)
}
