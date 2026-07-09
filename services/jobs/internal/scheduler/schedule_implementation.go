// Derived from https://github.com/roylee0704/gron (MIT License)
// Copyright (c) 2015 Roy Lee <roylee0704@gmail.com>
package scheduler

import "time"

type schedulePeriodic struct {
	period time.Duration
}

func (schedule schedulePeriodic) Next(t time.Time) time.Time {
	return t.Truncate(time.Second).Add(schedule.period)
}

func (schedule schedulePeriodic) At(t string) Schedule {
	if schedule.period < time.Hour*24 {
		panic("period must be at least in days")
	}

	tp, err := time.Parse("15:04", t)
	if err != nil {
		panic("invalid time format, expected hh:mm")
	}

	return &scheduleAt{
		period: schedule.period,
		hh:     tp.Hour(),
		mm:     tp.Minute(),
	}
}

type scheduleAt struct {
	period time.Duration
	hh     int
	mm     int
}

func (schedule scheduleAt) reset(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), schedule.hh, schedule.mm, 0, 0, t.Location())
}

func (schedule scheduleAt) Next(t time.Time) time.Time {
	next := schedule.reset(t)
	if !t.Before(next) {
		return next.Add(schedule.period)
	}
	return next
}

type scheduleOnce struct {
	delay time.Duration
	run   bool
}

func (schedule *scheduleOnce) Next(t time.Time) time.Time {
	if schedule.run {
		return time.Time{}
	}
	schedule.run = true
	return t.Add(schedule.delay)
}
