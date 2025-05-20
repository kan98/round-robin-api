package config

import (
	"os"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

var once sync.Once
var config *Config

type Config struct {
	ApiPorts         []string
	LoadBalancerPort string
	OptimiseConnPool bool
}

func Get() *Config {
	once.Do(func() {
		godotenv.Load()

		config = &Config{
			ApiPorts:         getPorts(),
			LoadBalancerPort: os.Getenv("loadBalancerPort"),
			OptimiseConnPool: os.Getenv("optimiseConnPool") == "true",
		}
	})

	return config
}

func getPorts() []string {
	portsListRaw := os.Getenv("apiPorts")

	portsList := strings.Split(portsListRaw, ",")

	result := []string{}
	for _, str := range portsList {
		if str != "" {
			result = append(result, str)
		}
	}
	return result
}

func Reset() {
	once = sync.Once{}
	config = nil
}
