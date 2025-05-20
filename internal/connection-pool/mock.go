package connectionpool

import (
	"errors"
	"net/http"
	"time"
)

type MockConnectionPool struct {
	GetConnectionToReturnErr bool
}

type MockConnection struct {
}

func (m *MockConnectionPool) GetConnection() (Connection, error) {
	if m.GetConnectionToReturnErr {
		return nil, errors.New("error")
	}
	return &MockConnection{}, nil
}

func (m *MockConnection) GetUrl() string {
	return "http://localhost:2222"
}

func (m *MockConnection) Analyse(resp *http.Response, err error, startTime time.Time) {
}
