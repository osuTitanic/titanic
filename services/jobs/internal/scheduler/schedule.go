// Derived from https://github.com/roylee0704/gron (MIT License)
// Copyright (c) 2015 Roy Lee <roylee0704@gmail.com>
package scheduler

import (
	"time"
)

// Schedule is the interface that wraps the basic Next method.
type Schedule interface {
	// Next calculates the next occurrence of the schedule based on the provided time t.
	// It returns a zero time (time.Time{}) if the schedule should never run again.
	Next(t time.Time) time.Time
}

// AtSchedule extends Schedule by enabling periodic-interval & time-specific setup.
type AtSchedule interface {
	// At returns a schedule that reoccurs at the specified time of day, e.g. "15:04".
	At(t string) Schedule
	Schedule
}

// Once returns a Schedule that runs once after the given duration.
func Once(in time.Duration) Schedule {
	return &scheduleOnce{delay: in}
}

// Now returns a Schedule that runs once immediately.
func Now() Schedule {
	return Once(0)
}

// Every returns a Schedule reoccurs every given period.
func Every(period time.Duration) AtSchedule {
	if period < time.Second {
		period = time.Second
	}
	period = period - time.Duration(period.Nanoseconds())%time.Second
	return &schedulePeriodic{period: period}
}
