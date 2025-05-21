package seedsimulator

import (
	"math/rand"

	"kan.com/round-robin-api/internal/config"
)

type Seed struct {
	seedNo             int64
	ProbabilityToErr   float64
	AverageSleepSpeed  int
	ProbabilityOfSleep float64
}

func New(index int) *Seed {
	seedCfg := config.Get().ApiSimulatorSeed
	if seedCfg == 0 {
		return nil
	}

	seedNo := seedCfg + int64(index)
	r := rand.New(rand.NewSource(seedNo))

	// Randomise err rate between 0 and 0.5
	probabilityToErr := r.Float64() / 2

	// Randomise latency between 0 and 500ms
	averageSleepSpeed := r.Intn(500)

	// Randomise probability of sleep between 0 and 0.5
	probabilityOfSleep := r.Float64() / 2

	return &Seed{
		ProbabilityToErr:   probabilityToErr,
		AverageSleepSpeed:  averageSleepSpeed,
		ProbabilityOfSleep: probabilityOfSleep,
		seedNo:             seedNo,
	}
}

func (s *Seed) ToError() bool {
	if s == nil {
		return false
	}

	randomFloat := rand.Float64()

	return randomFloat < s.ProbabilityToErr
}

func (s *Seed) SleepTime() int {
	if s != nil {
		randomFloat := rand.Float64()
		if randomFloat < s.ProbabilityOfSleep {
			return s.AverageSleepSpeed
		}
	}

	return 0
}
