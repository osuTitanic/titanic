package replays

import (
	"math"
	"sort"
)

const (
	minimumTouchscreenPresses = 100
	jumpDistanceThreshold     = 80.0
	jumpSpeedThreshold        = 5.0
	pressWindowMilliseconds   = 50
)

const gameplayButtons = Left1 | Right1 | Left2 | Right2

// DetectTouchscreenUsage analyzes replay frames for touchscreen-like cursor
// movement and returns whether the resulting score meets decisionThreshold.
func DetectTouchscreenUsage(frames []Frame, decisionThreshold float64) (detected bool, score float64) {
	speeds := make([]float64, 0, max(len(frames)-1, 0))
	teleportCount := 0
	pressCount := 0
	pressesAfterTeleport := 0
	lastTeleportTime := 0
	hasTeleport := false

	// We check for 3 main signals for analysis:
	// 1. Large fast cursor jumps, i.e. a movement that is considered a "teleport"
	// 2. Presses shortly after a teleport, a sign of touchscreen usage
	// 3. Very high cursor speeds, distributed across the replay (checking the 95th percentile here)

	// This method of analysis can still produce false positives, but I think its a good starting point
	// especially when only checking for high decision thresholds (0.8 by default)

	for index := 1; index < len(frames); index++ {
		previous := frames[index-1]
		current := frames[index]

		// A movement sample represents a single movement between two
		// frames, which we can analyze for speed and distance
		speed, distance, usable := calculateMovementSample(previous, current)
		if usable {
			speeds = append(speeds, speed)

			if isTeleportMovement(speed, distance) {
				teleportCount++
				lastTeleportTime = current.Time
				hasTeleport = true
			}
		}

		if isNewButtonPress(previous.Buttons, current.Buttons) {
			pressCount++

			if isPressAfterTeleport(current.Time, lastTeleportTime, hasTeleport) {
				pressesAfterTeleport++
			}
		}
	}

	if len(speeds) == 0 {
		// No usable movement samples were found, likely due to a malformed replay
		return false, 0
	}

	if pressCount < minimumTouchscreenPresses {
		// Not enough button presses for replay analysis to be reliable
		return false, 0
	}

	teleportRatio := float64(teleportCount) / float64(len(speeds))
	pressTeleportRatio := float64(pressesAfterTeleport) / float64(pressCount)
	p95Speed := calculatePercentile(speeds, 95)

	score = calculateTouchscreenScore(teleportRatio, pressTeleportRatio, p95Speed)
	return score >= decisionThreshold, score
}

func calculateMovementSample(previous, current Frame) (speed, distance float64, usable bool) {
	delta := current.Time - previous.Time
	if delta <= 0 {
		return 0, 0, false
	}

	// Distance can be calculated using our good old friend pythagoras
	// Then we just use some simple 5th grade physics to calculate speed = distance / time
	distance = math.Hypot(current.X-previous.X, current.Y-previous.Y)
	speed = distance / float64(delta)
	return speed, distance, true
}

func isTeleportMovement(speed, distance float64) bool {
	return distance >= jumpDistanceThreshold && speed >= jumpSpeedThreshold
}

func isNewButtonPress(previous, current ButtonState) bool {
	previousGameplay := previous & gameplayButtons
	currentGameplay := current & gameplayButtons
	return currentGameplay&^previousGameplay != 0
}

func isPressAfterTeleport(pressTime, lastTeleportTime int, hasTeleport bool) bool {
	if !hasTeleport {
		return false
	}
	timeSinceTeleport := pressTime - lastTeleportTime
	return timeSinceTeleport >= 0 && timeSinceTeleport <= pressWindowMilliseconds
}

func calculateTouchscreenScore(teleportRatio, pressTeleportRatio, p95Speed float64) float64 {
	teleportScore := min(teleportRatio/0.06, 1.0)
	pressScore := min(pressTeleportRatio/0.30, 1.0)
	speedScore := min(p95Speed/4.0, 1.0)

	// TODO: figure out weighting for these factors, currently just a guess
	return teleportScore*0.35 + pressScore*0.45 + speedScore*0.20
}

// Adapted from https://github.com/montanaflynn/stats/blob/master/percentile.go (MIT)
// Copyright (c) 2014-2026 Montana Flynn (https://montanaflynn.com)
func calculatePercentile(input []float64, percent float64) float64 {
	length := len(input)
	if length == 0 {
		return 0
	}
	if length == 1 {
		return input[0]
	}

	// Start by sorting a copy of the slice
	c := append([]float64(nil), input...)
	sort.Float64s(c)

	// Use linear interpolation between the closest ranks
	rank := (percent / 100) * float64(length-1)
	k := int(rank)
	f := rank - float64(k)
	if k+1 < length {
		return c[k] + f*(c[k+1]-c[k])
	}

	return c[k]
}
