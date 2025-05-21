package config

import (
	"os"
	"strconv"
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
	ApiSimulatorSeed int64
}

func Get() *Config {
	once.Do(func() {
		godotenv.Load()

		config = &Config{
			ApiPorts:         getPorts(),
			LoadBalancerPort: os.Getenv("loadBalancerPort"),
			OptimiseConnPool: os.Getenv("optimiseConnPool") == "true",
			ApiSimulatorSeed: getApiSimulatorSeed(),
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

func getApiSimulatorSeed() int64 {
	seed := os.Getenv("apiSimulatorSeed")
	if seed == "" {
		return 0
	}

	seedInt, err := strconv.ParseInt(seed, 10, 64)
	if err != nil {
		return 0
	}
	return seedInt
}

func Reset() {
	once = sync.Once{}
	config = nil
}
