package seedsimulator

import (
	"math/rand"

	"kan.com/round-robin-api/internal/config"
)

type Seed struct {
	ProbabilityToErr   float64
	AverageSleepSpeed  int
	ProbabilityOfSleep float64
}

var seedNo int64 = -1

func New() *Seed {
	if seedNo == -1 {
		seedNo = config.Get().ApiSimulatorSeed
	} else {
		seedNo++
	}

	if seedNo == 0 {
		return nil
	}

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
