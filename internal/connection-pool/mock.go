package connectionpool

import (
	"net/http"
	"time"
)

type MockConnectionPool struct {
}

type MockConnection struct {
}

func (m *MockConnectionPool) GetConnection() (Connection, error) {
	return &MockConnection{}, nil
}

func (m *MockConnection) GetUrl() string {
	return "http://localhost:2222"
}

func (m *MockConnection) Analyse(resp *http.Response, err error, startTime time.Time) {
}
