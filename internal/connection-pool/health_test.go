package connectionpool

import (
	"testing"
)

func TestHealth(t *testing.T) {
	runAsync = false
	globalLatencySum = 0
	globalLatencyCount = 0

	t.Run("updates the score when error", func(t *testing.T) {
		h := health{}
		h.analyse(true, 100)

		if h.score != 10 {
			t.Errorf("expected score to be 10, got %d", h.score)
		}
	})

	t.Run("updates the score when no error", func(t *testing.T) {
		h := health{}
		h.analyse(false, 100)

		if h.score != -10 {
			t.Errorf("expected score to be -10, got %d", h.score)
		}
	})

	t.Run("updates the score when latency is lower than avg", func(t *testing.T) {
		globalLatencySum = 1000
		globalLatencyCount = 10

		h := health{}

		h.analyse(false, 50)

		if h.score > 0 {
			t.Errorf("expected score to be more than 0, got %d", h.score)
		}
	})

	t.Run("updates score on latency is higher than avg", func(t *testing.T) {
		globalLatencySum = 1000
		globalLatencyCount = 10

		h := health{}
		h.analyse(false, 800)

		if h.score < 0 {
			t.Errorf("expected score to be less than 0, got %d", h.score)
		}
	})

	t.Run("adds penalty when score is over threshold", func(t *testing.T) {
		h := health{score: scorePenaltyThreshold - 1}

		h.analyse(true, 25)

		if h.score != 0 {
			t.Errorf("score should reset to 0 but got %d", h.score)
		}
		if h.penaltyLeftTimes <= 0 {
			t.Errorf("penalty should be more than 0 but got %d", h.penaltyLeftTimes)
		}
	})

	t.Run("score cannot go below the min threshold", func(t *testing.T) {
		h := health{score: -3000}

		h.analyse(false, 0)

		if h.score != -20 {
			t.Errorf("score is lower than the minimum threshold %d", h.score)
		}
	})

	t.Run("penalty cannot go above the max threshold", func(t *testing.T) {
		h := health{score: 3000}

		h.analyse(false, 0)

		if h.score != 0 {
			t.Errorf("score should reset to 0 but got %d", h.score)
		}

		if h.penaltyLeftTimes != 50 {
			t.Errorf("penalty is higher than the maximum threshold %d", h.penaltyLeftTimes)
		}
	})

	t.Run("if conn still has penalty, we should decrement penalty by 1", func(t *testing.T) {
		h := health{score: 0, penaltyLeftTimes: 50}

		h.analyse(false, 0)

		if h.score != 0 {
			t.Errorf("score should stay at 0 but got %d", h.score)
		}

		if h.penaltyLeftTimes != 49 {
			t.Errorf("penalty should be 49 but is %d", h.penaltyLeftTimes)
		}
	})
}

func TestIsHealthy(t *testing.T) {
	t.Run("healthy when penalty is 0", func(t *testing.T) {
		h := health{penaltyLeftTimes: 0}
		if !h.isHealthy() {
			t.Errorf("should be healthy")
		}
	})

	t.Run("not healthy when penalty is more than 0", func(t *testing.T) {
		h := health{penaltyLeftTimes: 1}
		if h.isHealthy() {
			t.Errorf("should not be healthy")
		}
	})
}

func TestDecreasePenalty(t *testing.T) {
	t.Run("decreases penalty when penalty is more than 0", func(t *testing.T) {
		h := health{penaltyLeftTimes: 2}
		h.decreasePenalty()

		if h.penaltyLeftTimes != 1 {
			t.Errorf("penalty should be decreased to 1 but got %d", h.penaltyLeftTimes)
		}
	})

	t.Run("does not decrease penalty when penalty is 0", func(t *testing.T) {
		h := health{penaltyLeftTimes: 0}
		h.decreasePenalty()

		if h.penaltyLeftTimes != 0 {
			t.Errorf("penalty should be 0 but got %d", h.penaltyLeftTimes)
		}
	})
}
