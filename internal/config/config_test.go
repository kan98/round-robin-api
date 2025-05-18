package config

import (
	"reflect"
	"sync"
	"testing"
)

func TestConfig(t *testing.T) {
	t.Run("Valid env variables", func(t *testing.T) {
		reset()
		t.Setenv("apiPorts", "8080,8081,8082")
		t.Setenv("loadBalancerPort", "8083")
		t.Setenv("optimiseConnPool", "true")

		config := Get()

		if !reflect.DeepEqual(config.ApiPorts, []string{"8080", "8081", "8082"}) {
			t.Error("Incorrect api ports value")
		}

		if config.LoadBalancerPort != "8083" {
			t.Error("incorrect load balancer port value")
		}

		if !config.OptimiseConnPool {
			t.Error("incorrect optimiseConnPool value")
		}
	})

	t.Run("Empty env variables", func(t *testing.T) {
		reset()
		t.Setenv("apiPorts", "")
		t.Setenv("loadBalancerPort", "")
		t.Setenv("optimiseConnPool", "")

		config := Get()

		if !reflect.DeepEqual(config.ApiPorts, []string{}) {
			t.Error("Incorrect api ports value")
		}

		if config.LoadBalancerPort != "" {
			t.Error("incorrect load balancer port value")
		}

		if config.OptimiseConnPool {
			t.Error("incorrect optimiseConnPool value")
		}
	})

	t.Run("apiPorts trims empty ports in list correctly", func(t *testing.T) {
		reset()
		t.Setenv("apiPorts", "8080,,,8081,8082")

		config := Get()

		if !reflect.DeepEqual(config.ApiPorts, []string{"8080", "8081", "8082"}) {
			t.Error("Incorrect api ports value")
		}
	})
}

func reset() {
	once = sync.Once{}
	config = nil
}
