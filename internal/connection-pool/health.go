package connectionpool

import (
	"sync"
)

const (
	maxGlobalLatencyCount = 10
	scorePenaltyThreshold = 20
	minAllowedScore       = -20
	maxAllowedScore       = 50
)

var (
	globalLatencyMutex = &sync.Mutex{}
	globalLatencySum   int64
	globalLatencyCount int64
	runAsync           = true
)

type health struct {
	score            int
	penaltyLeftTimes int
}

func (h *health) isHealthy() bool {
	return h.penaltyLeftTimes <= 0
}

func (h *health) decreasePenalty() {
	if h.penaltyLeftTimes > 0 {
		h.penaltyLeftTimes--
	}
}

func (h *health) analyse(isErr bool, latencyMs int64) {
	analyseFunc := func() {
		if h.penaltyLeftTimes > 0 {
			h.penaltyLeftTimes--
			return
		}

		if isErr {
			h.score += 10
		} else {
			h.score -= 10
		}

		avg := avgGlobalLatency(latencyMs)
		if avg > 0 {
			latencyDiffPercent := (latencyMs - avg) * 100 / avg
			h.score += int(latencyDiffPercent / 10)
		}

		if h.score < minAllowedScore {
			h.score = minAllowedScore
		} else if h.score > maxAllowedScore {
			h.score = maxAllowedScore
		}

		if h.score >= scorePenaltyThreshold {
			h.penaltyLeftTimes = h.score
			h.score = 0
		}
	}

	// to run unit tests synchronously
	if runAsync {
		go analyseFunc()
	} else {
		analyseFunc()
	}
}

func avgGlobalLatency(currentLatencyMs int64) int64 {
	globalLatencyMutex.Lock()
	defer globalLatencyMutex.Unlock()

	if globalLatencyCount == 0 {
		addNewLatency(currentLatencyMs)
		return currentLatencyMs
	}

	avg := globalLatencySum / globalLatencyCount
	addNewLatency(currentLatencyMs)
	return avg
}

func addNewLatency(currentLatencyMs int64) {
	if globalLatencyCount >= maxGlobalLatencyCount {
		globalLatencySum = (globalLatencySum*9)/10 + currentLatencyMs
	} else {
		globalLatencySum += currentLatencyMs
		globalLatencyCount++
	}
}
