package connectionpool

import (
	"errors"
	"net/url"
)

type connectionPool struct {
	connections  []connection
	currentIndex int
}

func New(ports []string) *connectionPool {
	connections := make([]connection, 0)

	for _, port := range ports {
		connections = append(connections, connection{
			url: &url.URL{
				Scheme: scheme,
				Host:   hostname + port,
			},
		})
	}

	return &connectionPool{
		connections:  connections,
		currentIndex: 0,
	}
}

func (cp *connectionPool) GetConnection() (Connection, error) {
	if len(cp.connections) == 0 {
		return nil, errors.New("no available connections")
	}

	startIndex := cp.currentIndex

	for i := 0; i < len(cp.connections); i++ {
		index := (startIndex + i) % len(cp.connections)

		if cp.connections[index].health.isHealthy() {
			cp.currentIndex = (index + 1) % len(cp.connections)
			return &cp.connections[index], nil
		} else {
			cp.connections[index].health.decreasePenalty()
		}
	}

	// If all conns unhealthy, just return the next one
	index := cp.currentIndex
	cp.currentIndex = (cp.currentIndex + 1) % len(cp.connections)
	return &cp.connections[index], nil
}
