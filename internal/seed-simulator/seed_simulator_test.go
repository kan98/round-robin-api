package seedsimulator

import (
	"testing"

	"kan.com/round-robin-api/internal/config"
)

func TestNew(t *testing.T) {
	t.Run("seed is valid", func(t *testing.T) {
		t.Setenv("apiSimulatorSeed", "12345")
		config.Reset()

		seed := New(0)
		if seed == nil {
			t.Error("seed should not be nil")
		} else {
			if seed.ProbabilityToErr < 0 || seed.ProbabilityToErr > 0.5 {
				t.Errorf("probabilityToErr should be between 0 and 0.5 but is %f", seed.ProbabilityToErr)
			}
			if seed.AverageSleepSpeed < 0 || seed.AverageSleepSpeed > 500 {
				t.Errorf("averageSleepSpeed should be between 0 and 500 but is %d", seed.AverageSleepSpeed)
			}
			if seed.ProbabilityOfSleep < 0 || seed.ProbabilityOfSleep > 0.5 {
				t.Errorf("probabilityOfSleep should be between 0 and 0.5 but is %f", seed.ProbabilityOfSleep)
			}
		}
	})
	t.Run("second seed in runtime is incremented", func(t *testing.T) {
		t.Setenv("apiSimulatorSeed", "12345")
		config.Reset()

		seed1 := New(0)
		if seed1.seedNo != 12345 {
			t.Errorf("seedNo should be 12345 but is %d", seed1.seedNo)
		}

		seed2 := New(1)
		if seed2.seedNo != 12346 {
			t.Errorf("seedNo should be 12346 but is %d", seed2.seedNo)
		}
	})
	t.Run("seed is 0", func(t *testing.T) {
		t.Setenv("apiSimulatorSeed", "0")
		config.Reset()

		seed := New(0)
		if seed != nil {
			t.Error("seed should be nil")
		}
	})
}

func TestToError(t *testing.T) {
	t.Run("seed is nil", func(t *testing.T) {
		t.Setenv("apiSimulatorSeed", "0")
		config.Reset()

		seed := New(0)
		if seed.ToError() {
			t.Error("should return false for nil seed")
		}
	})
	t.Run("seed is valid", func(t *testing.T) {
		seed := Seed{ProbabilityToErr: 1}

		if !seed.ToError() {
			t.Error("should return false for valid seed with probability 1")
		}
	})
}
func TestSleepTime(t *testing.T) {
	t.Run("seed is nil", func(t *testing.T) {
		t.Setenv("apiSimulatorSeed", "0")
		config.Reset()

		seed := New(0)
		if seed.SleepTime() != 0 {
			t.Error("should return 0 for nil seed")
		}
	})
	t.Run("seed is valid", func(t *testing.T) {
		seed := Seed{AverageSleepSpeed: 100, ProbabilityOfSleep: 1}

		if seed.SleepTime() != 100 {
			t.Error("should return 100 for valid seed with probability 1")
		}
	})
}
