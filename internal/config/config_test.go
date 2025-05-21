package config

import (
	"reflect"
	"testing"
)

func TestConfig(t *testing.T) {
	t.Run("Valid env variables", func(t *testing.T) {
		Reset()
		t.Setenv("apiPorts", "2222,3333,4444")
		t.Setenv("loadBalancerPort", "8083")
		t.Setenv("optimiseConnPool", "true")
		t.Setenv("apiSimulatorSeed", "123")

		config := Get()

		if !reflect.DeepEqual(config.ApiPorts, []string{"2222", "3333", "4444"}) {
			t.Error("Incorrect api ports value")
		}

		if config.LoadBalancerPort != "8083" {
			t.Error("incorrect load balancer port value")
		}

		if !config.OptimiseConnPool {
			t.Error("incorrect optimiseConnPool value")
		}

		if config.ApiSimulatorSeed != 123 {
			t.Error("incorrect apiSimulatorSeed value")
		}
	})

	t.Run("Empty env variables", func(t *testing.T) {
		Reset()
		t.Setenv("apiPorts", "")
		t.Setenv("loadBalancerPort", "")
		t.Setenv("optimiseConnPool", "")
		t.Setenv("apiSimulatorSeed", "")

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

		if config.ApiSimulatorSeed != 0 {
			t.Error("incorrect apiSimulatorSeed value")
		}
	})

	t.Run("apiSimulatorSeed returns 0 for non number value", func(t *testing.T) {
		Reset()
		t.Setenv("apiSimulatorSeed", "abc")

		config := Get()

		if config.ApiSimulatorSeed != 0 {
			t.Error("incorrect apiSimulatorSeed value")
		}
	})

	t.Run("apiPorts trims empty ports in list correctly", func(t *testing.T) {
		Reset()
		t.Setenv("apiPorts", "2222,,,3333,4444")

		config := Get()

		if !reflect.DeepEqual(config.ApiPorts, []string{"2222", "3333", "4444"}) {
			t.Error("Incorrect api ports value")
		}
	})
}

func TestReset(t *testing.T) {
	t.Run("Reset config", func(t *testing.T) {
		Reset()
		t.Setenv("loadBalancerPort", "2222")
		port := Get().LoadBalancerPort

		if port != "2222" {
			t.Error("incorrect load balancer port value")
		}

		Reset()
		t.Setenv("loadBalancerPort", "3333")
		if Get().LoadBalancerPort != "3333" {
			t.Error("incorrect load balancer port value")
		}
	})
}
